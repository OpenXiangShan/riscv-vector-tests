package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeRdVs2Vm(pos int) []string {
	combinations := i.combinations([]LMUL{1}, []SEW{8}, []bool{false, true})

	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		vs2 := int(c.LMUL1) * 2
		builder.WriteString(i.gWriteRandomData(LMUL(3)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, c.LMUL1, SEW(8)))

		builder.WriteString(fmt.Sprintf("addi a0, a0, %d\n", 1*i.vlenb()))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(fmt.Sprintf("addi a0, a0, %d\n", 1*i.vlenb()))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, SEW(8)))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s s0, v%d%s\n",
			i.Name, vs2, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gMoveScalarToVector("s0", vd, SEW(i.Option.XLEN)))

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, SEW(i.Option.XLEN)))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}
	return res
}
