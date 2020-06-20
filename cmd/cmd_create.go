/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hitzhangjie/gorpc/config"
	"github.com/hitzhangjie/gorpc/descriptor"
	"github.com/hitzhangjie/gorpc/params"
	"github.com/hitzhangjie/gorpc/parser"
	"github.com/hitzhangjie/gorpc/tpl"
	"github.com/hitzhangjie/gorpc/util/fs"
	"github.com/hitzhangjie/gorpc/util/lang"
	"github.com/hitzhangjie/gorpc/util/log"
	"github.com/hitzhangjie/gorpc/util/pb"
	"github.com/hitzhangjie/gorpc/util/style"
	"github.com/hitzhangjie/gorpc/util/swagger"
)

var (
	createOption    *params.Option
	createSuccess   bool
	createOutputDir string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "指定pb文件快速创建工程或rpcstub",
	Long: `指定pb文件快速创建工程或rpcstub，

'gorpc create' 有两种模式:
- 生成一个完整的服务工程
- 生成被调服务的rpcstub，需指定'-rpconly'选项.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// 解析命令行参数
		err := cmd.ParseFlags(args)
		if err != nil {
			return fmt.Errorf("parse flags error: %v", err)
		}

		// 检查命令行参数
		option, err := loadCreateOption(cmd.Flags())
		if err != nil {
			return fmt.Errorf("check flags error: %v", err)
		}
		createOption = option

		// 初始化日志级别
		log.InitLogging(option.Verbose)
		log.Info("ready to process protofile: %s", option.ProtofileAbs)

		// 解析pb
		fd, err := parser.ParseProtoFile(option)
		if err != nil {
			return fmt.Errorf("Parse protofile: %s error: %v", option.Protofile, err)
		}
		fd.FilePath = option.ProtofileAbs
		fd.Dump()

		// 创建工程
		var outputdir string
		if !option.RpcOnly {
			outputdir, err = create(fd, option)
		} else {
			outputdir, err = generateRPCStub(fd, option)
		}

		if err != nil {
			if !option.RpcOnly {
				return fmt.Errorf("create gorpc project error: %v", err)
			}
			return fmt.Errorf("create gorpc stub error: %v", err)
		}

		// 生成 swagger 文档
		if option.SwaggerOn {
			if err = swagger.GenSwagger(fd, option); err != nil {
				return fmt.Errorf("create swagger api document error: %v", err)
			}
		}

		createOption = option
		createOutputDir = outputdir
		createSuccess = true

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {

		if !createSuccess {
			return nil
		}

		if createOption.Language != "go" {
			return nil
		}

		err := os.Chdir(createOutputDir)
		if err != nil {
			return err
		}

		// run mockgen to generate interface mock stub
		// fixme if `go build` not run first, `go generate` which calls `mockgen`,
		// it's really time consuming, so forget it so far.
		//
		//_, err = exec.LookPath("mockgen")
		//if err != nil {
		//	log.Error("please install mockgen in order to generate mockstub")
		//} else {
		//	c := exec.Command("go", "generate")
		//	if buf, err := c.CombinedOutput(); err != nil {
		//		log.Error("run mockgen error: %v,\n%s", err, string(buf))
		//		return err
		//	}
		//}

		// run goimports to format your code
		goimports, err := exec.LookPath("goimports")
		if err != nil {
			log.Error("please install goimports in order to format code")
		} else {
			c := exec.Command(goimports, "-w", ".")
			if buf, err := c.CombinedOutput(); err != nil {
				log.Error("run goimports error: %v,\n%s", err, string(buf))
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringArray("protodir", []string{"."}, "include path of the target protofile")
	createCmd.Flags().StringP("protofile", "p", "", "protofile used as IDL of target service")
	createCmd.Flags().String("protocol", "gorpc", "protocol to use, gorpc, http, etc")
	createCmd.Flags().BoolP("verbose", "v", false, "show verbose logging info")
	createCmd.Flags().String("assetdir", "", "path of project template")
	createCmd.Flags().Bool("alias", false, "enable alias mode of rpc name")
	createCmd.Flags().Bool("rpconly", false, "generate rpc stub only")
	createCmd.Flags().String("lang", "go", "programming language, including go, java, python")
	createCmd.Flags().StringP("mod", "m", "", "go module, default: ${pb.package}")
	createCmd.Flags().StringP("output", "o", "", "output directory")
	createCmd.Flags().BoolP("force", "f", false, "enable overwritten existed code forcibly")
	createCmd.Flags().Bool("swagger", false, "enable swagger to gen swagger api document.")
}

func loadCreateOption(flagSet *pflag.FlagSet) (*params.Option, error) {

	option := loadCreateFlagsetToOption(flagSet)

	// 检查pb文件是否合法
	if len(option.Protofile) == 0 {
		return nil, errors.New("invalid protofile")
	}

	// 定位pb文件
	target, err := fs.LocateFile(option.Protofile, option.Protodirs)
	if err != nil {
		return nil, err
	}
	option.Protofile = filepath.Base(target)
	option.ProtofileAbs = target
	option.Protodirs = append(option.Protodirs, filepath.Dir(target))

	// 加载gorpc.json中定义的语言相关的配置
	option.GoRPCCfg, err = config.GetLanguageCfg(option.Language)
	if err != nil {
		return nil, fmt.Errorf("load config via gorpc.json error: %v", err)
	}
	if len(option.Assetdir) == 0 {
		option.Assetdir = option.GoRPCCfg.AssetDir
	}

	// 判断gomod
	// - 优先使用-mod指定的moduleName
	// - 没有指定-mod选项的话，再考虑加载本地go.mod，兼容老的操作逻辑
	// - 如果本地也没有指定go.mod，再考虑pb中的package（模板实现的）
	if len(option.GoMod) == 0 {
		mod, err := lang.LoadGoMod()
		if err == nil && len(mod) != 0 {
			option.GoModEx = mod
			option.GoMod = mod
		}
	}

	return option, nil
}

func loadCreateFlagsetToOption(flagSet *pflag.FlagSet) *params.Option {

	option := &params.Option{}

	option.Protodirs, _ = flagSet.GetStringArray("protodir")
	option.Protofile, _ = flagSet.GetString("protofile")
	option.Language, _ = flagSet.GetString("lang")
	option.Protocol, _ = flagSet.GetString("protocol")
	option.AliasOn, _ = flagSet.GetBool("alias")
	option.RpcOnly, _ = flagSet.GetBool("rpconly")
	option.Assetdir, _ = flagSet.GetString("assetdir")
	option.Verbose, _ = flagSet.GetBool("verbose")
	option.GoMod, _ = flagSet.GetString("mod")
	option.OutputDir, _ = flagSet.GetString("output")
	option.Force, _ = flagSet.GetBool("force")
	option.SwaggerOn, _ = flagSet.GetBool("swagger")

	return option
}

// create 代码生成，生成完整的工程
func create(fd *descriptor.FileDescriptor, option *params.Option) (outputdir string, err error) {

	// - 准备输出目录
	outputdir, err = getOutputDir(option)
	if err != nil {
		return
	}

	if !isSafeOutputDir(outputdir) && !option.Force {
		err = fmt.Errorf("reject overwrite existed code: %s,\nuse --force/-f to make it if you want", outputdir)
		return
	}

	defer func() {
		if err != nil {
			removeDirAsNeeded(outputdir)
		}
	}()

	// - 生成代码
	err = tpl.GenerateFiles(fd, option, outputdir)
	if err != nil {
		return
	}

	// create rpcstub
	stubDir := filepath.Join(outputdir, "stub")
	if _, err = os.Lstat(stubDir); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(stubDir, os.ModePerm); err != nil {
			return
		}
	}
	stub := filepath.Join(outputdir, "stub")

	// - move outputdir/rpc to outputdir/stub/dir($gopkgdir)
	fileOption := fmt.Sprintf("%s_package", option.GoRPCCfg.Language)
	pbPackage, err := parser.GetPbPackage(fd, fileOption)
	if err != nil {
		return
	}

	if fileOption == "java_package" {
		pathLast := path.Join(strings.Split(pbPackage, ".")...)
		pbPackage = path.Join("client/src/main/java", strings.ToLower(pathLast))
	} else if fileOption == "python_package" {
		pbPackage = strings.Replace(pbPackage, ".", "_", -1)
	}

	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	pbOutDir := filepath.Join(stub, pbPackage)
	err = os.MkdirAll(pbOutDir, os.ModePerm)
	if err != nil {
		return
	}

	pb2pkg := fd.Pb2ImportPath

	// 检查pb中是否启用validation特性，并生成go文件
	resultCheckSECV := parser.CheckSECVEnabled(fd)
	if resultCheckSECV {
		err = pb.SECVProtoc(option.Protodirs, option.Protofile, option.Language, pbOutDir, pb2pkg)
		if err != nil {
			removeDirAsNeeded(outputdir)
			err = fmt.Errorf("GenerateFiles: %v", err)
			return
		}
	}

	// 处理-protofile指定的pb文件
	err = pb.Protoc(fd, option.Protodirs, option.Protofile, option.Language, pbOutDir, pb2pkg)
	if err != nil {
		err = fmt.Errorf("GenerateFiles: %v", err)
		return
	}

	// - copy *.proto to outpoutdir/rpc/
	basename := filepath.Base(fd.FilePath)
	err = fs.Copy(fd.FilePath, filepath.Join(pbOutDir, basename))
	if err != nil {
		return
	}

	// - 处理${protofile}依赖的其他pb文件
	//BUG: 目录组织问题，不再按照pb相对路径关系进行组织，全部按照stub/package进行组织
	//err = handleDependencies(fd, option, pbPackage, pbOutDir)
	err = handleDependencies(fd, option, pbPackage, stub)
	if err != nil {
		return
	}

	// fixme 用抽象、泛化来解决特殊逻辑问题，这里为java调整pb路径的逻辑需要调整下，go的类似
	err = changeProtofileDir(pbOutDir, option.GoRPCCfg.Language)
	if err != nil {
		return
	}

	// - 将outputdir/rpc移动到outputdir/$gopkgdir/
	src := filepath.Join(outputdir, "rpc")
	defer os.RemoveAll(src)
	dest := path.Join(stub, pbPackage)

	// 找到生成的validate.pb.go文件，并移动到stub目录下
	if resultCheckSECV {

		err = filepath.Walk(pbPackage, func(fpath string, info os.FileInfo, err error) (e error) {

			protoName := strings.Split(option.Protofile, ".")[0]
			generatedName := protoName + ".pb.validate.go"
			srcSameDirPath := "./" + outputdir + ".pb.validate.go"
			_, checkSameDir := os.Stat(srcSameDirPath)

			if checkSameDir != nil {

				if fname := filepath.Base(fpath); fname == generatedName {
					destDir := filepath.Join(dest, filepath.Base(fpath))
					moveResult := fs.Move(fpath, destDir)

					// - 生成结束后，删除validate所在的目录
					parentDirArr := strings.Split(fpath, "/")
					parentDir := parentDirArr[0]

					err = os.RemoveAll(parentDir)

					return moveResult
				}

				return nil

			}

			moveResult := fs.Move(srcSameDirPath, dest)

			// - 生成结束后，删除validate所在的目录
			parentDirArr := strings.Split(fpath, "/")
			parentDir := parentDirArr[0]

			err = os.RemoveAll(parentDir)

			return moveResult

		})

		if err != nil {
			return
		}
	}

	err = filepath.Walk(src, func(fpath string, info os.FileInfo, err error) (e error) {

		if fpath == src {
			return nil
		}

		if fname := filepath.Base(fpath); fname == "gorpc.go" {
			// - 将stub文件gorpc.go重命名，
			fname = fs.BaseNameWithoutExt(basename)
			return fs.Move(fpath, filepath.Join(dest, fname+".gorpc.go"))
		} else {
			return fs.Move(fpath, filepath.Join(dest, filepath.Base(fpath)))
		}
	})
	if err != nil {
		return
	}

	// Python 移动stub的setup.py
	if option.GoRPCCfg.Language == "python" {
		// move stub/.../setup.py to stub/setup.py
		setupFilePath := filepath.Join(pbOutDir, "setup.py")

		var fin os.FileInfo
		if fin, err = os.Stat(setupFilePath); err == nil && !fin.IsDir() {
			log.Debug("move setup.py file from %s to %s", setupFilePath, stubDir)
			if err = fs.Move(setupFilePath, stubDir); err != nil {
				return
			}
		}
	}

	// Java格式化操作
	if option.GoRPCCfg.Language == "java" {
		if err = javaFormat(outputdir, fd, option); err != nil {
			return
		}
	}

	// 格式化操作
	if err = style.GoFmtDir(outputdir); err != nil {
		return
	}

	log.Info("Generate project %s```%s```%s success", log.COLOR_RED, basename, log.COLOR_GREEN)
	return
}

func generateRPCStub(fd *descriptor.FileDescriptor, option *params.Option) (outputdir string, err error) {

	// 代码生成
	// - 准备输出目录
	outputdir, err = os.Getwd()
	if err != nil {
		return
	}

	if filepath.IsAbs(option.OutputDir) {
		outputdir = option.OutputDir
	} else {
		outputdir = filepath.Join(outputdir, option.OutputDir)
	}
	err = os.MkdirAll(outputdir, os.ModePerm)
	if err != nil {
		return
	}

	// - 生成代码，只处理clientstub
	generated := map[string]struct{}{}
	for _, f := range option.GoRPCCfg.RPCClientStub {
		in := filepath.Join(option.Assetdir, f)
		log.Debug("handle:%s", in)
		out := filepath.Join(outputdir, strings.TrimSuffix(filepath.Base(in), option.GoRPCCfg.TplFileExt))
		if err = tpl.GenerateFile(fd, in, out, option); err != nil {
			return
		}

		if err = style.Format(out, option); err != nil {
			return
		}

		generated[out] = struct{}{}
	}

	// 将stub文件gorpc.go重命名
	basename := fs.BaseNameWithoutExt(fd.FilePath)

	// fix issue: https://github.com/hitzhangjie/gorpc/issues/178
	//
	// 通过generated记录生成的文件列表，解决上述issue，有可能会尝试去重命名path/to/go/pkg/mod/.../gorpc.go，
	// 这是错误的，只应该遍历生成的文件列表并对生成的gorpc.go进行重命名。
	// 另外，filepath.Walk的方式可能遍历了不该遍历的文件，比如通过软连接指向pb文件父目录的情况下，然后在HOME下
	// 执行gorpc create命令，由于默认的outputdir为当前目录（-o可修改输出目录），此时会遍历HOME下文件，非常耗时。
	//
	// 这里直接遍历生成的文件列表即可。
	//err = filepath.Walk(outputdir, func(fpath string, info os.FileInfo, err error) (e error) {
	//	if _, ok := generated[fpath]; !ok {
	//		return
	//	}
	//	if fname := filepath.Base(fpath); fname == "gorpc.go" {
	//		e = fs.Move(fpath, filepath.Join(path.Dir(fpath), basename+".gorpc.go"))
	//	}
	//	return
	//})
	//if err != nil {
	//	return err
	//}

	for fpath, _ := range generated {
		if fname := filepath.Base(fpath); fname != "gorpc.go" {
			continue
		}
		dst := filepath.Join(path.Dir(fpath), basename+".gorpc.go")
		if err = fs.Move(fpath, dst); err != nil {
			return
		}
		break
	}

	//将所有package相同的依赖过滤出来
	var protofiles []string
	protofiles = append(protofiles, option.Protofile)
	for fname, pkg := range fd.Pb2ValidGoPkg {
		if pkg == fd.PackageName {
			protofiles = append(protofiles, fname)
		}
	}
	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	if err = pb.Protoc(fd, option.Protodirs, option.Protofile, option.Language, outputdir, fd.Pb2ImportPath); err != nil {
		//if err = pb.Protoc(c.Option.Protodirs, c.Option.Protofile, c.Option.Language, outputdir, fd.Pb2ValidGoPkg); err != nil {
		err = fmt.Errorf("GenerateFiles: %v", err)
		return
	}
	log.Info("Generate rpc stub success")
	return
}

func isSafeOutputDir(dir string) bool {

	_, err := os.Lstat(dir)

	// 目录不存在，说明不存在写覆盖的情况
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		return false
	}

	// 存在的话，检测下目录下是否存在源文件，存在就有覆盖风险
	err = filepath.Walk(dir, func(p string, inf os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".go") {
			return fmt.Errorf("go code detected: %s", p)
		}
		return nil
	})

	if err != nil {
		return false
	}
	return true
}

func removeDirAsNeeded(path string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if dir == path {
		return
	}
	os.RemoveAll(path)
}

// handleDependencies 处理-protofile指定的pb文件中import的其他pb文件，包括protoc处理，与拷贝pb文件
//
// 准备protoc生成pb文件对应的*.pb.go，需要注意的是，避免生成循环依赖的代码
//
// jhump/protoreflect解析结果，如果是与-protofile相同的pkgname，那么importpath为"",
//
// protoc --go_out=M$pb=$pkgname，这里需要做兼容处理:
// 	1. 避免传递$pkgname为空, 否则protoc会生成这样的代码：
//    ```go
//    package $pkgname
//    import (
//        "."
//    )
//    ```
// 	2. 避免传递与-protofile相同pkgname的情况，不然会导致循环依赖:
//    ```go
//    package $pkgname
//    import (
//        "$pkgname"
//    )
//    ```
func handleDependencies(fd *descriptor.FileDescriptor, option *params.Option, pbPackage string, outputDir string) (err error) {
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return err
	}

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return err
	}

	includeDirs := []string{}
	for fname, _ := range fd.Pb2ImportPath {
		dir, _ := filepath.Split(fname)
		includeDirs = append(includeDirs, dir)
	}

	// 计算所依赖文件的package
	//dependPackage := map[string]string{}
	//for k, p := range fd.Pb2ImportPath {
	//	dependPackage[k] = p
	//}
	for fname, importPath := range fd.Pb2ImportPath {

		// 如果是${protofile}跳过不处理
		if fname == fd.FilePath {
			continue
		}

		// 跳过google官方提供的pb文件，gorpc扩展文件，swagger 扩展文件
		if strings.HasPrefix(fname, "google/protobuf") ||
			fname == "gorpc.proto" || fname == "validate.proto" || fname == "swagger.proto" {
			continue
		}

		pbOutDir := filepath.Join(outputDir, importPath)
		if option.Language == "java" {
			pbOutDir = filepath.Join(outputDir, pbPackage)
		}
		if err := os.MkdirAll(pbOutDir, os.ModePerm); err != nil {
			return err
		}

		// 继承上一级的目录,避免出现目录找不到的问题
		searchPath := option.Protodirs

		parentDirs := []string{wd}
		parentDirs = append(parentDirs, option.Protodirs...)

		for _, pDir := range parentDirs {
			for _, incDir := range includeDirs {

				includeDir := filepath.Join(pDir, incDir)
				includeDir = filepath.Clean(includeDir)

				if fin, err := os.Lstat(includeDir); err != nil {
					if !os.IsNotExist(err) {
						return err
					}
				} else {
					if !fin.IsDir() {
						return fmt.Errorf("import path: %s, not directory", includeDir)
					}
					searchPath = append(searchPath, includeDir)
				}
			}
		}

		//if err := pb.Protoc(searchPath, fname, option.Language, pbOutDir, pb2ValidGoPkg); err != nil {
		//	return fmt.Errorf("GenerateFiles: %v", err)
		//}
		if err := pb.Protoc(fd, searchPath, fname, option.Language, pbOutDir, fd.Pb2ImportPath); err != nil {
			return fmt.Errorf("GenerateFiles: %v", err)
		}

		// 拷贝pb文件
		p, err := fs.LocateFile(fname, option.Protodirs)
		if err != nil {
			return err
		}

		_, baseName := filepath.Split(fname)
		src := p
		dst := filepath.Join(pbOutDir, baseName)

		log.Debug("Copy file %s to %s", src, dst)
		if err := fs.Copy(src, dst); err != nil {
			return err
		}

		// 初始化gomod
		//
		// 避免重复初始化go.mod
		fp := filepath.Join(pbOutDir, "go.mod")
		fin, err := os.Stat(fp)
		if err == nil && !fin.IsDir() {
			continue
		}

		// fixme 移动到createCmd.PostRun
		// 执行go mod init, 与pbPackage相同也不用初始化
		if option.Language != "go" {
			continue
		}

		if len(importPath) != 0 && importPath != pbPackage {
			os.Chdir(pbOutDir)

			cmd := exec.Command("go", "mod", "init", importPath)
			if buf, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("process %s, initialize go.mod in stub/%s error: %v", fname, importPath, string(buf))
			}
			log.Debug("process %s, initialize go.mod success in xxxout/%s: go mod init %s", fname, importPath, importPath)
		}
	}

	if err = os.Chdir(wd); err != nil {
		return err
	}

	return nil
}

type DirMove struct {
	Src string
	Dst string
}

// java生成代码格式化处理，当前先hardcode，后续考虑通过脚本配置实现
func javaFormat(outputdir string, fd *descriptor.FileDescriptor, option *params.Option) (err error) {
	log.Debug("******************************java format begin***********************************")
	serviceName := fd.Services[0].Name
	packageName := serviceName
	if javaPackage, ok := fd.FileOptions["java_package"].(string); ok {
		ss := strings.Split(javaPackage, ".")
		packageName = filepath.Join(ss...) // java路径调整 suggest from youngwwang
		//idx := len(ss) - 1
		//packageName = ss[idx]
	}
	var DirMoveList []DirMove
	err = filepath.Walk(outputdir, func(fpath string, info os.FileInfo, err error) (e error) {
		if info.IsDir() && strings.HasSuffix(fpath, "gorpcserver") {
			dstPath := strings.TrimRight(fpath, "gorpcserver") + strings.ToLower(packageName)
			DirMoveList = append(DirMoveList, DirMove{fpath, dstPath})
		} else if !info.IsDir() && strings.HasSuffix(fpath, "service_api.java") {
			dstPath := strings.TrimRight(fpath, "service_api.java") + strcase.ToCamel(serviceName) + "Api." + option.GoRPCCfg.Language
			e = os.RemoveAll(dstPath)
			log.Debug("file move, src: %s, dst: %s", fpath, dstPath)
			e = fs.Move(fpath, dstPath)
		}
		return
	})
	for _, dirMove := range DirMoveList {
		log.Debug("dir copy, src: %s, dst: %s", dirMove.Src, dirMove.Dst)
		err = fs.Copy(dirMove.Src, dirMove.Dst)
		err = os.RemoveAll(dirMove.Src)
	}
	log.Debug("******************************java format finish***********************************")
	return
}