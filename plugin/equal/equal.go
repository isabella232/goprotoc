// Copyright (c) 2013, Vastech SA (PTY) LTD. All rights reserved.
// http://code.google.com/p/gogoprotobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
The equal plugin generates an Equal and a VerboseEqual method for each message.
These equal methods are quite obvious.
The only difference is that VerboseEqual returns a non nil error if it is not equal.
This error contains more detail on exactly which part of the message was not equal to the other message.
The idea is that this is useful for debugging.

Equal is enabled using the following extensions:

  - equal
  - equal_all

While VerboseEqual is enable dusing the following extensions:

  - verbose_equal
  - verbose_equal_all

The equal plugin also generates a test given it is enabled using one of the following extensions:

  - testgen
  - testgen_all

Let us look at:

  code.google.com/p/gogoprotobuf/test/example/example.proto

Btw all the output can be seen at:

  code.google.com/p/gogoprotobuf/test/example/*

The following message:

  option (gogoproto.equal_all) = true;
  option (gogoproto.verbose_equal_all) = true;

  message B {
	option (gogoproto.description) = true;
	optional string A = 1 [(gogoproto.embed) = true];
	repeated int64 G = 2 [(gogoproto.customtype) = "github.com/dropbox/goprotoc/test.Id"];
  }

given to the equal plugin, will generate the following code:

	func (this *B) VerboseEqual(that interface{}) error {
		if that == nil {
			if this == nil {
				return nil
			}
			return fmt1.Errorf("that == nil && this != nil")
		}

		that1, ok := that.(*B)
		if !ok {
			return fmt1.Errorf("that is not of type *B")
		}
		if that1 == nil {
			if this == nil {
				return nil
			}
			return fmt1.Errorf("that is type *B but is nil && this != nil")
		} else if this == nil {
			return fmt1.Errorf("that is type *Bbut is not nil && this == nil")
		}
		if this.xxx_IsASet != that1.xxx_IsASet {
			return fmt1.Errorf("that.a is not equal to this.a")
		}
		if this.xxx_IsASet && this.a != that1.a {
			return fmt1.Errorf("a this(%v) Not Equal that(%v)", this.a, that1.a)
		}
		if this.xxx_LenG != that1.xxx_LenG {
			return fmt1.Errorf("that.g is not equal to this.g")
		}
		for i := 0; i < this.xxx_LenG; i++ {
			if this.g[i] != that1.g[i] {
				return fmt1.Errorf("g this[%v](%v) Not Equal that[%v](%v)", i, this.g[i], i, that1.g[i])
			}
		}
		if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
			return fmt1.Errorf("XXX_unrecognized this(%v) Not Equal that(%v)", this.XXX_unrecognized, that1.XXX_unrecognized)
		}
		return nil
	}

	func (this *B) Equal(that interface{}) bool {
		if that == nil {
			if this == nil {
				return true
			}
			return false
		}

		that1, ok := that.(*B)
		if !ok {
			return false
		}
		if that1 == nil {
			if this == nil {
				return true
			}
			return false
		} else if this == nil {
			return false
		}
		if this.xxx_IsASet != that1.xxx_IsASet {
			return false
		}
		if this.xxx_IsASet && this.a != that1.a {
			return false
		}
		if this.xxx_LenG != that1.xxx_LenG {
			return false
		}
		for i := 0; i < this.xxx_LenG; i++ {
			if this.g[i] != that1.g[i] {
				return false
			}
		}
		if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
			return false
		}
		return true
	}

and the following test code:

	func TestBVerboseEqual(t *testing8.T) {
		popr := math_rand8.New(math_rand8.NewSource(time8.Now().UnixNano()))
		p := NewPopulatedB(popr, false)
		data, err := dropbox_gogoprotobuf_proto2.Marshal(p)
		if err != nil {
			panic(err)
		}
		msg := &B{}
		if err := dropbox_gogoprotobuf_proto2.Unmarshal(data, msg); err != nil {
			panic(err)
		}
		if err := p.VerboseEqual(msg); err != nil {
			t.Fatalf("%#v !VerboseEqual %#v, since %v", msg, p, err)
	}

*/
package equal

import (
	"github.com/dropbox/goprotoc/gogoproto"
	"github.com/dropbox/goprotoc/protoc-gen-dgo/generator"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	bytesPkg generator.Single
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return "equal"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.bytesPkg = p.NewImport("bytes")

	for _, msg := range file.Messages() {
		if gogoproto.HasVerboseEqual(file.FileDescriptorProto, msg.DescriptorProto) {
			p.generateMessage(msg, true, gogoproto.HasExtensionsMap(file.FileDescriptorProto, msg.DescriptorProto))
		}
		if gogoproto.HasEqual(file.FileDescriptorProto, msg.DescriptorProto) {
			p.generateMessage(msg, false, gogoproto.HasExtensionsMap(file.FileDescriptorProto, msg.DescriptorProto))
		}
	}
}

func (p *plugin) generateMessage(message *generator.Descriptor, verbose bool, hasExtensionsMap bool) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())
	if verbose {
		p.P(`func (this *`, ccTypeName, `) VerboseEqual(that interface{}) error {`)
	} else {
		p.P(`func (this *`, ccTypeName, `) Equal(that interface{}) bool {`)
	}
	p.In()
	p.P(`if that == nil {`)
	p.In()
	p.P(`if this == nil {`)
	p.In()
	if verbose {
		p.P(`return nil`)
	} else {
		p.P(`return true`)
	}
	p.Out()
	p.P(`}`)
	if verbose {
		p.P(`return `, p.Pkg["fmt"], `.Errorf("that == nil && this != nil")`)
	} else {
		p.P(`return false`)
	}
	p.Out()
	p.P(`}`)
	p.P(``)
	p.P(`that1, ok := that.(*`, ccTypeName, `)`)
	p.P(`if !ok {`)
	p.In()
	if verbose {
		p.P(`return `, p.Pkg["fmt"], `.Errorf("that is not of type *`, ccTypeName, `")`)
	} else {
		p.P(`return false`)
	}
	p.Out()
	p.P(`}`)
	p.P(`if that1 == nil {`)
	p.In()
	p.P(`if this == nil {`)
	p.In()
	if verbose {
		p.P(`return nil`)
	} else {
		p.P(`return true`)
	}
	p.Out()
	p.P(`}`)
	if verbose {
		p.P(`return `, p.Pkg["fmt"], `.Errorf("that is type *`, ccTypeName, ` but is nil && this != nil")`)
	} else {
		p.P(`return false`)
	}
	p.Out()
	p.P(`} else if this == nil {`)
	p.In()
	if verbose {
		p.P(`return `, p.Pkg["fmt"], `.Errorf("that is type *`, ccTypeName, `but is not nil && this == nil")`)
	} else {
		p.P(`return false`)
	}
	p.Out()
	p.P(`}`)

	for _, field := range message.Field {
		fieldname := p.GetFieldName(message, field)
		repeated := field.IsRepeated()
		if repeated {
			p.P(`if this.`, generator.SizerName(fieldname), ` != that1.`, generator.SizerName(fieldname), ` {`)
		} else {
			p.P(`if this.`, generator.SetterName(fieldname), ` != that1.`, generator.SetterName(fieldname), ` {`)
		}
		p.In()
		if verbose {
			p.P(`return `, p.Pkg["fmt"], `.Errorf("that.`, fieldname, ` is not equal to this.`, fieldname, `")`)
		} else {
			p.P(`return false`)
		}
		p.Out()
		p.P(`}`)

		if !repeated {
			if field.IsMessage() || p.IsGroup(field) {
				p.P(`if this.`, generator.SetterName(fieldname), ` && !this.`, fieldname, `.Equal(that1.`, fieldname, `) {`)
			} else if field.IsBytes() {
				p.P(`if this.`, generator.SetterName(fieldname), ` && !`, p.bytesPkg.Use(), `.Equal(this.`, fieldname, `, that1.`, fieldname, `) {`)
			} else {
				p.P(`if this.`, generator.SetterName(fieldname), ` && this.`, fieldname, ` != that1.`, fieldname, `{`)
			}
			p.In()
			if verbose {
				p.P(`return `, p.Pkg["fmt"], `.Errorf("`, fieldname, ` this(%v) Not Equal that(%v)", this.`, fieldname, `, that1.`, fieldname, `)`)
			} else {
				p.P(`return false`)
			}
			p.Out()
			p.P(`}`)
		} else {
			p.P(`for i := 0; i < this.`, generator.SizerName(fieldname), `; i++ {`)
			p.In()
			if field.IsMessage() || p.IsGroup(field) {
				p.P(`if !this.`, fieldname, `[i].Equal(that1.`, fieldname, `[i]) {`)
			} else if field.IsBytes() {
				p.P(`if !`, p.bytesPkg.Use(), `.Equal(this.`, fieldname, `[i], that1.`, fieldname, `[i]) {`)
			} else {
				p.P(`if this.`, fieldname, `[i] != that1.`, fieldname, `[i] {`)
			}
			p.In()
			if verbose {
				p.P(`return `, p.Pkg["fmt"], `.Errorf("`, fieldname, ` this[%v](%v) Not Equal that[%v](%v)", i, this.`, fieldname, `[i], i, that1.`, fieldname, `[i])`)
			} else {
				p.P(`return false`)
			}
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
		}
	}
	if message.DescriptorProto.HasExtension() {
		fieldname := "XXX_extensions"
		if hasExtensionsMap {
			p.P(`for k, v := range this.`, fieldname, ` {`)
			p.In()
			p.P(`if v2, ok := that1.`, fieldname, `[k]; ok {`)
			p.In()
			p.P(`if !v.Equal(&v2) {`)
			p.In()
			if verbose {
				p.P(`return `, p.Pkg["fmt"], `.Errorf("`, fieldname, ` this[%v](%v) Not Equal that[%v](%v)", k, this.`, fieldname, `[k], k, that1.`, fieldname, `[k])`)
			} else {
				p.P(`return false`)
			}
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`} else  {`)
			p.In()
			if verbose {
				p.P(`return `, p.Pkg["fmt"], `.Errorf("`, fieldname, `[%v] Not In that", k)`)
			} else {
				p.P(`return false`)
			}
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)

			p.P(`for k, _ := range that1.`, fieldname, ` {`)
			p.In()
			p.P(`if _, ok := this.`, fieldname, `[k]; !ok {`)
			p.In()
			if verbose {
				p.P(`return `, p.Pkg["fmt"], `.Errorf("`, fieldname, `[%v] Not In this", k)`)
			} else {
				p.P(`return false`)
			}
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
		} else {
			p.P(`if !`, p.bytesPkg.Use(), `.Equal(this.`, fieldname, `, that1.`, fieldname, `) {`)
			p.In()
			if verbose {
				p.P(`return `, p.Pkg["fmt"], `.Errorf("`, fieldname, ` this(%v) Not Equal that(%v)", this.`, fieldname, `, that1.`, fieldname, `)`)
			} else {
				p.P(`return false`)
			}
			p.Out()
			p.P(`}`)
		}
	}
	fieldname := "XXX_unrecognized"
	p.P(`if !`, p.bytesPkg.Use(), `.Equal(this.`, fieldname, `, that1.`, fieldname, `) {`)
	p.In()
	if verbose {
		p.P(`return `, p.Pkg["fmt"], `.Errorf("`, fieldname, ` this(%v) Not Equal that(%v)", this.`, fieldname, `, that1.`, fieldname, `)`)
	} else {
		p.P(`return false`)
	}
	p.Out()
	p.P(`}`)
	if verbose {
		p.P(`return nil`)
	} else {
		p.P(`return true`)
	}
	p.Out()
	p.P(`}`)
}

func init() {
	generator.RegisterPlugin(NewPlugin())
}
