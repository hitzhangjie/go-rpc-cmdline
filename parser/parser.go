package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hitzhangjie/codeblocks/log"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/extension/gorpc"
	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/util/lang"
	"github.com/hitzhangjie/gorpc-cli/util/pb"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

// ParseProtoFile 解析proto文件，返回一个构造好的可以应用于模板填充的FileDescriptor对象
//
// ParseProtoFile负责的工作包括：
// - 解析pb文件，拿到原始的描述信息
// - 检查工程约束，如是否制定了go_option、method option等自定义的一些业务开发约束
func ParseProtoFile(option *params.Option) (*descriptor.FileDescriptor, error) {

	protodirs := option.Protodirs

	p, err := pb.LocateGoRPCProto()
	if err != nil {
		return nil, err
	}
	protodirs = append(protodirs, p)

	// 解析pb
	var fd *desc.FileDescriptor
	if fds, err := parseProtoFile(option.Protofile, protodirs...); err != nil {
		return nil, err
	} else {
		fd = fds[0]
	}
	// 检查约束
	if err := checkRequirements(fd); err != nil {
		return nil, err
	}

	// 构造可以用于指导代码生成的FileDescriptor
	fileDescriptor := new(descriptor.FileDescriptor)
	// 设置依赖(import的pb文件及其输出包名)
	fillDependencies(fd, fileDescriptor)
	// - 设置packageName
	withErrorCheck(fillPackageName(fd, fileDescriptor))
	// - 设置imports
	withErrorCheck(fillImports(fd, fileDescriptor))
	// - 设置fileOptions
	withErrorCheck(fillFileOptions(fd, fileDescriptor))
	// - 设置service
	withErrorCheck(fillServices(fd, fileDescriptor))
	// - 设置app server
	withErrorCheck(fillAppServerName(fd, fileDescriptor))
	// - 设置rpc请求响应类型对应的定义pb
	withErrorCheck(fillRPCMessageTypes(fd, fileDescriptor))

	fileDescriptor.SetRawFileDescriptor(fd)

	return fileDescriptor, nil
}

// parseProtoFile 调用jhump/protoreflect来解析pb文件，拿到proto文件描述信息
func parseProtoFile(fname string, protodirs ...string) ([]*desc.FileDescriptor, error) {

	parser := protoparse.Parser{
		ImportPaths:           protodirs,
		IncludeSourceCodeInfo: true,
	}

	return parser.ParseFiles(fname)
}

// checkRequirements 检查是否符合某些约束条件
//
// requirements:
// - pb语法必须为syntax = "proto3";
// - 必须指定fileoption go_package;
// - service定义数量不能为0;
func checkRequirements(fd *desc.FileDescriptor) error {

	// fixme MUST: syntax = "proto3"
	//if !fd.IsProto3() {
	//	return errors.New("syntax isn't proto3")
	//}

	// fixme MUST: option go_package = "github.com/$group/$repo"
	// fixme MUST: option go_package trailing part, must equal to package directive
	//opts := fd.GetFileOptions()
	//if opts == nil {
	//	return errors.New(`FileOption 'go_package' missing`)
	//}
	//
	//gopkg := opts.GetGoPackage()
	//if len(gopkg) == 0 {
	//	return errors.New(`FileOption 'go_package' missing`)
	//} else {
	//	var trailing string
	//	idx := strings.LastIndex(gopkg, "/")
	//
	//	if idx < 0 {
	//		trailing = gopkg
	//	} else {
	//		trailing = gopkg[idx+1:]
	//	}
	//
	//	if trailing != fd.GetPackage() {
	//		return errors.New(`'option go_package="a/b/c"' trailing part "c" must be consistent with 'package diretive'`)
	//	}
	//}

	// MUST: service
	if len(fd.GetServices()) == 0 {
		return errors.New("service missing")
	}

	return nil
}

func withErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func fillDependencies(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) error {

	pb2ValidGoPkg := map[string]string{} // k=pb文件名, v=protoc处理后package名
	pb2ImportPath := map[string]string{} // k=pb文件名, v=go代码中对应importpath

	func() {
		validGoPkg := lang.PBValidGoPackage(fd.GetPackage())
		importPath := fd.GetPackage()

		if opts := fd.GetFileOptions(); opts != nil {
			if gopkgopt := opts.GetGoPackage(); len(gopkgopt) != 0 {
				validGoPkg = lang.PBValidGoPackage(gopkgopt)
				importPath = gopkgopt
			}
		}
		pb2ValidGoPkg[fd.GetName()] = validGoPkg
		pb2ImportPath[fd.GetName()] = importPath
	}()

	var f func(*desc.FileDescriptor)

	f = func(fd *desc.FileDescriptor) {

		for _, dep := range fd.GetDependencies() {
			if len(dep.GetDependencies()) != 0 {
				f(dep)
			}

			fname := dep.GetFullyQualifiedName()
			pkg := dep.GetPackage()

			pb2ValidGoPkg[fname] = lang.PBValidGoPackage(pkg)

			var (
				validGoPkg = lang.PBValidGoPackage(pkg)
				importPath = pkg
			)
			if opts := dep.GetFileOptions(); opts != nil {
				if gopkgopt := opts.GetGoPackage(); len(gopkgopt) != 0 {
					validGoPkg = lang.PBValidGoPackage(gopkgopt)
					importPath = gopkgopt
				}
			}
			pb2ValidGoPkg[fname] = validGoPkg
			pb2ImportPath[fname] = importPath
		}
	}

	f(fd)

	nfd.Pb2ValidGoPkg = pb2ValidGoPkg
	nfd.Pb2ImportPath = pb2ImportPath

	return nil
}

func fillPackageName(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) error {
	nfd.PackageName = fd.GetPackage()
	return nil
}

func fillAppServerName(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) error {
	strs := strings.Split(fd.GetPackage(), ".")
	if len(strs) == 3 {
		nfd.AppName = strs[1]
		nfd.ServerName = strs[2]
	}
	return nil
}

func fillImports(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) error {
	nfd.Imports = getImports(fd, nfd)
	return nil
}

func fillFileOptions(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) error {

	opts := fd.GetFileOptions()
	if opts == nil {
		return nil
	}

	v, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(v, &m); err != nil {
		return err
	}

	if nfd.FileOptions == nil {
		nfd.FileOptions = make(map[string]interface{})
	}

	for k, v := range m {
		nfd.FileOptions[k] = v
	}
	return nil
}

func fillServices(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) error {

	for _, sd := range fd.GetServices() {

		nsd := new(descriptor.ServiceDescriptor)
		nfd.Services = append(nfd.Services, nsd)

		// service name
		nsd.Name = sd.GetName()

		// service methods
		for _, m := range sd.GetMethods() {

			leadingComments := strings.TrimSpace(m.GetSourceInfo().GetLeadingComments())
			trailingComments := strings.TrimSpace(m.GetSourceInfo().GetTrailingComments())

			rpc := &descriptor.RPCDescriptor{
				Name: m.GetName(),
				Cmd:  m.GetName(),
				// fixme 这里写死了rpc的拼接规则为/$package.$service/$method
				FullyQualifiedCmd: fmt.Sprintf("/%s.%s/%s", fd.GetPackage(), sd.GetName(), m.GetName()),
				RequestType:       m.GetInputType().GetFullyQualifiedName(),
				ResponseType:      m.GetOutputType().GetFullyQualifiedName(),
				//go1.12+
				//LeadingComments:   strings.ReplaceAll(leadingComments, "\n", "\n// "),
				//TrailingComments:  strings.ReplaceAll(trailingComments, "\n", "\n// "),
				//
				//compatible with go1.11
				LeadingComments:  strings.Replace(leadingComments, "\n", "\n// ", -1),
				TrailingComments: strings.Replace(trailingComments, "\n", "\n// ", -1),
			}
			nsd.RPC = append(nsd.RPC, rpc)

			// check method option, if gorpc.alias exists, use it as rpc.baseCmd
			if v, err := proto.GetExtension(m.GetMethodOptions(), gorpc.E_Alias); err == nil {
				s := v.(*string)
				if s == nil {
					log.Debug("method:%s.%s parse methodOptions option gorpc.alias not specified", sd.GetName(), m.GetName())
				} else {
					log.Debug("method:%s.%s parse methodOptions, name:%v = %s", sd.GetName(), m.GetName(), gorpc.E_Alias, *(v.(*string)))
					if s != nil {
						if cmd := strings.TrimSpace(*s); len(cmd) != 0 {
							rpc.FullyQualifiedCmd = cmd
						}
					}
				}
			}

			// check method option, if gorpc.swagger exists, fill the rpc.swagger_info
			if v, err := proto.GetExtension(m.GetMethodOptions(), gorpc.E_Swagger); err == nil {
				swagger := v.(*gorpc.SwaggerRule)
				if swagger == nil {
					log.Debug("method:%s.%s parse methodOptions option gorpc.swagger not specified", sd.GetName(), m.GetName())
				} else {
					if title := strings.TrimSpace(swagger.Title); len(title) == 0 {
						// 如果title 为空，这里会取 rpc 定义的前注释作为方法的 title
						rpc.SwaggerInfo.Title = strings.Replace(leadingComments, "\n", "\n// ", -1)
					}
					rpc.SwaggerInfo.Title = strings.TrimSpace(swagger.Title)
					rpc.SwaggerInfo.Description = strings.TrimSpace(swagger.Description)
					if method := strings.TrimSpace(swagger.Method); len(method) == 0 {
						// FIXME 如果 method 为空的话，为了支持 swagger-ui 显示，必须定义一个方法。默认 POST
						rpc.SwaggerInfo.Method = "post"
					} else {
						rpc.SwaggerInfo.Method = method
					}
					log.Debug("method:%s.%s parse methodOptions, title: %s, method: %s, description: %s. ",
						sd.GetName(), m.GetName(), rpc.SwaggerInfo.Title, rpc.SwaggerInfo.Method, rpc.SwaggerInfo.Description)
				}
			} else {
				// 如果没有定义 rpc 的 swagger option，则需要从注释中填充
				rpc.SwaggerInfo.Title = strings.Replace(leadingComments, "\n", "\n// ", -1)
				// 默认为 post
				rpc.SwaggerInfo.Method = "post"
				rpc.SwaggerInfo.Description = ""
			}
		}
	}

	return nil
}

// fillRPCMessageTypes 桩代码里面涉及到rpc请求、响应类型名，这个要建立与定义他们的pb的映射关系
func fillRPCMessageTypes(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) error {

	def := map[string]string{}

	for _, sd := range fd.GetServices() {
		for _, m := range sd.GetMethods() {
			in := m.GetInputType().GetFullyQualifiedName()
			out := m.GetOutputType().GetFullyQualifiedName()

			inDefLoc, err := findMessageDefLocation(in, fd)
			if err != nil {
				return err
			}
			def[in] = inDefLoc

			outDefLoc, err := findMessageDefLocation(out, fd)
			if err != nil {
				return err
			}
			def[out] = outDefLoc
		}
	}

	if len(def) != 0 {
		nfd.RpcMessageType = def
	}
	return nil
}

func findMessageDefLocation(typ string, fd *desc.FileDescriptor) (string, error) {
	for _, t := range fd.GetMessageTypes() {
		if t.GetFullyQualifiedName() == typ {
			return fd.GetFullyQualifiedName(), nil
		}
	}

	for _, dep := range fd.GetDependencies() {
		for _, t := range dep.GetMessageTypes() {
			if t.GetFullyQualifiedName() == typ {
				return dep.GetFullyQualifiedName(), nil
			}
		}
	}

	return "", errors.New("not found")
}

func getImports(fd *desc.FileDescriptor, nfd *descriptor.FileDescriptor) []string {
	imports := []string{}

	// 避免针对同一个package重复多次import，goimports可以去掉`import but unused`问题，
	// 但是解决不了`redeclared as imported package name`这类问题
	existed := map[string]struct{}{}

	for _, dep := range fd.GetDependencies() {

		pb := dep.GetName()
		pbImport, ok := nfd.Pb2ImportPath[pb]
		if !ok {
			panic(fmt.Errorf("get import path of %s fail", pb))
		}

		_, ok = existed[pbImport]
		if !ok {
			imports = append(imports, pbImport)
			existed[pbImport] = struct{}{}
		}
	}

	return imports
}

// GetPbPackage 获取pb放置的路径
func GetPbPackage(fd *descriptor.FileDescriptor, fileOption string) (string, error) {

	pbPackage := fd.PackageName
	if fo := fd.FileOptions[fileOption]; fo != nil {
		if v := fd.FileOptions[fileOption].(string); len(v) != 0 {
			pbPackage = v
		}
	}

	return pbPackage, nil
}
