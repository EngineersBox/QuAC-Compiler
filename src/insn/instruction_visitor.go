package insn

import (
	"fmt"
	"strconv"

	"github.com/EngineersBox/QuAC-Compiler/antlr4"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type InsnVisitor struct {
	antlr4.QuACParserVisitor
	labels map[string]uint16
}

func NewInsnVisitor(labels map[string]uint16) InsnVisitor {
	return InsnVisitor{
		labels: labels,
	}
}

func (this *InsnVisitor) Visit(tree antlr.ParseTree) interface{} {
	switch val := tree.(type) {
	case *antlr4.ParseContext:
		return this.VisitParse(val)
	case *antlr4.IFormatStatementContext:
		return this.VisitIFormatStatement(val)
	case *antlr4.RMemFormatStatementContext:
		return this.VisitRMemFormatStatement(val)
	case *antlr4.RALUFormatStatementContext:
		return this.VisitRALUFormatStatement(val)
	case *antlr4.NopStatementContext:
		return this.VisitNopStatement(val)
	case *antlr4.Pseudo2ParamStatementContext:
		return this.VisitPseudo2ParamStatement(val)
	case *antlr4.JprStatementContext:
		return this.VisitJprStatement(val)
	case *antlr4.JpmStatementContext:
		return this.VisitJpmStatement(val)
	case *antlr4.JpStatementContext:
		return this.VisitJpStatement(val)
	case *antlr4.WordStatementContext:
		return this.VisitWordStatement(val)
	case *antlr4.LabelStatementContext:
		return this.VisitLabelStatement(val)
	case *antlr4.IFormatContext:
		return this.VisitIFormat(val)
	case *antlr4.RMemFormatContext:
		return this.VisitRMemFormat(val)
	case *antlr4.RALUFormatContext:
		return this.VisitRALUFormat(val)
	case *antlr4.NopContext:
		return this.VisitNop(val)
	case *antlr4.Pseudo2ParamContext:
		return this.VisitPseudo2Param(val)
	case *antlr4.JprContext:
		return this.VisitJpr(val)
	case *antlr4.JpmContext:
		return this.VisitJpm(val)
	case *antlr4.JpContext:
		return this.VisitJp(val)
	case *antlr4.RegisterContext:
		return this.VisitRegister(val)
	default:
		panic("Unknown context")
	}
}

func (this *InsnVisitor) VisitParse(ctx *antlr4.ParseContext) interface{} {
	var result []uint16 = make([]uint16, 0)
	for _, statement := range ctx.AllStatement() {
		val := this.Visit(statement)
		fmt.Printf("%04x\n", val)
		result = append(result, val.([]uint16)...)
	}
	return result
}

func (this *InsnVisitor) VisitIFormatStatement(ctx *antlr4.IFormatStatementContext) interface{} {
	return this.Visit(ctx.IFormat())
}

func (this *InsnVisitor) VisitRMemFormatStatement(ctx *antlr4.RMemFormatStatementContext) interface{} {
	return this.Visit(ctx.RMemFormat())
}

func (this *InsnVisitor) VisitRALUFormatStatement(ctx *antlr4.RALUFormatStatementContext) interface{} {
	return this.Visit(ctx.RALUFormat())
}

func (this *InsnVisitor) VisitNopStatement(ctx *antlr4.NopStatementContext) interface{} {
	return this.Visit(ctx.Nop())
}

func (this *InsnVisitor) VisitPseudo2ParamStatement(ctx *antlr4.Pseudo2ParamStatementContext) interface{} {
	return this.Visit(ctx.Pseudo2Param())
}

func (this *InsnVisitor) VisitJprStatement(ctx *antlr4.JprStatementContext) interface{} {
	return this.Visit(ctx.Jpr())
}

func (this *InsnVisitor) VisitJpmStatement(ctx *antlr4.JpmStatementContext) interface{} {
	return this.Visit(ctx.Jpm())
}

func (this *InsnVisitor) VisitJpStatement(ctx *antlr4.JpStatementContext) interface{} {
	return this.Visit(ctx.Jp())
}

func (this *InsnVisitor) VisitWordStatement(ctx *antlr4.WordStatementContext) interface{} {
	val, err := strconv.ParseInt(ctx.IntegerLiteral().GetText(), 0, 16)
	if err != nil {
		panic(err)
	}
	return []uint16{uint16(val)}
}

func (this *InsnVisitor) VisitLabelStatement(ctx *antlr4.LabelStatementContext) interface{} {
	return make([]uint16, 0)
}

func (this *InsnVisitor) VisitIFormat(ctx *antlr4.IFormatContext) interface{} {
	var result uint16
	if ctx.MOVL() != nil {
		result = MOVL_MASK
	} else if ctx.SETH() != nil {
		result = SETH_MASK
	} else {
		panic("Invalid I-Format instruction")
	}
	if ctx.EQ() != nil {
		result |= TRUE_CONDITION
	}
	var rd uint16 = this.Visit(ctx.Register()).(uint16)
	result |= rd << 8
	imm8, err := parseImm8(ctx.IntegerLiteral().GetText())
	if err != nil {
		panic(fmt.Errorf("Bad immediate 8-bit format: %s", err))
	}
	result |= uint16(imm8)
	fmt.Printf("I-Format: %016b\n", result)
	return []uint16{result}
}

func (this *InsnVisitor) VisitRMemFormat(ctx *antlr4.RMemFormatContext) interface{} {
	var result uint16
	if ctx.STR() != nil {
		result = STR_MASK
	} else if ctx.LDR() != nil {
		result = LDR_MASK
	} else {
		panic("Invalid R-Format memory instruction")
	}
	if ctx.EQ() != nil {
		result |= TRUE_CONDITION
	}
	var registers []antlr4.IRegisterContext = ctx.AllRegister()
	if len(registers) != 2 {
		panic(fmt.Errorf("Expected 2 registers (rd, ra), got %d", len(registers)))
	}
	result |= this.Visit(registers[0]).(uint16) << RD_REGISTER_OFFSET
	result |= this.Visit(registers[1]).(uint16) << RA_REGISTER_OFFSET
	fmt.Printf("R-Format Mem: %016b\n", result)
	return []uint16{result}
}

func (this *InsnVisitor) VisitRALUFormat(ctx *antlr4.RALUFormatContext) interface{} {
	var result uint16
	if ctx.ADD() != nil {
		result = ADD_MASK
	} else if ctx.SUB() != nil {
		result = SUB_MASK
	} else if ctx.AND() != nil {
		result = AND_MASK
	} else if ctx.ORR() != nil {
		result = ORR_MASK
	} else {
		panic("Invalid R-Format ALU instruction")
	}
	if ctx.EQ() != nil {
		result |= TRUE_CONDITION
	}
	var registers []antlr4.IRegisterContext = ctx.AllRegister()
	if len(registers) != 3 {
		panic(fmt.Errorf("Expected 3 registers (rd, ra, rb), got %d", len(registers)))
	}
	result |= this.Visit(registers[0]).(uint16) << RD_REGISTER_OFFSET
	result |= this.Visit(registers[1]).(uint16) << RA_REGISTER_OFFSET
	result |= this.Visit(registers[2]).(uint16) << RB_REGISTER_OFFSET
	fmt.Printf("R-Format ALU: %016b\n", result)
	return []uint16{result}
}

func (v *InsnVisitor) VisitNop(ctx *antlr4.NopContext) interface{} {
	return make([]uint16, 0)
}

func (this *InsnVisitor) VisitPseudo2Param(ctx *antlr4.Pseudo2ParamContext) interface{} {
	var result uint16
	if ctx.MOV() != nil {
		result = MOV_MASK
	} else if ctx.CMP() != nil {
		result = CMP_MASK
	} else {
		panic("Invalid mov/cmp instruction")
	}
	if ctx.EQ() != nil {
		result |= TRUE_CONDITION
	}
	var registers []antlr4.IRegisterContext = ctx.AllRegister()
	if len(registers) != 2 {
		panic(fmt.Errorf("Expected 2 registers (rd, ra), got %d", len(registers)))
	}
	if ctx.MOV() != nil {
		result |= this.Visit(registers[0]).(uint16) << RD_REGISTER_OFFSET
		result |= this.Visit(registers[1]).(uint16) << RA_REGISTER_OFFSET
	} else if ctx.CMP() != nil {
		result |= this.Visit(registers[0]).(uint16) << RA_REGISTER_OFFSET
		result |= this.Visit(registers[1]).(uint16) << RB_REGISTER_OFFSET
	}
	return []uint16{result}
}

func (this *InsnVisitor) VisitJpr(ctx *antlr4.JprContext) interface{} {
	var result uint16 = JPR_MASK
	if ctx.EQ() != nil {
		result |= TRUE_CONDITION
	}
	result |= this.Visit(ctx.Register()).(uint16) << RA_REGISTER_OFFSET
	return []uint16{result}
}

func (this *InsnVisitor) VisitJpm(ctx *antlr4.JpmContext) interface{} {
	var result uint16 = JPM_MASK
	if ctx.EQ() != nil {
		result |= TRUE_CONDITION
	}
	result |= this.Visit(ctx.Register()).(uint16) << RA_REGISTER_OFFSET
	return []uint16{result}
}

func (this *InsnVisitor) VisitJp(ctx *antlr4.JpContext) interface{} {
	var result uint16 = JP_MASK
	if ctx.EQ() != nil {
		result |= TRUE_CONDITION
	}
	if ctx.IntegerLiteral() != nil {
		imm8, err := parseImm8(ctx.IntegerLiteral().GetText())
		if err != nil {
			panic(fmt.Errorf("Bad immediate 8-bit format: %s", err))
		}
		result |= uint16(imm8)
	} else if ctx.Identifier() != nil {
		offset, ok := this.labels[ctx.Identifier().GetText()]
		if !ok {
			panic(fmt.Errorf("No label declared with name \"%s\"", ctx.Identifier().GetText()))
		}
		result |= offset
	} else {
		panic("Invalid jp instruction, expected a register or immediate 8-bit value")
	}
	fmt.Printf("Jp %016b\n", result)
	return []uint16{result}
}

func (this *InsnVisitor) VisitRegister(ctx *antlr4.RegisterContext) interface{} {
	if ctx.RZ() != nil || ctx.R0() != nil {
		return uint16(0b000)
	} else if ctx.R1() != nil {
		return uint16(0b001)
	} else if ctx.R2() != nil {
		return uint16(0b010)
	} else if ctx.R3() != nil {
		return uint16(0b011)
	} else if ctx.R4() != nil {
		return uint16(0b100)
	} else if ctx.FL() != nil || ctx.R5() != nil {
		return uint16(0b101)
	} else if ctx.PC() != nil || ctx.R7() != nil {
		return uint16(0b111)
	}
	panic("Invalid register")
}

func parseImm8(arg string) (uint8, error) {
	imm8, err := strconv.ParseInt(arg, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(imm8), nil
}
