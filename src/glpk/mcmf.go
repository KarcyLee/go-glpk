package glpk

// #cgo LDFLAGS: -lglpk
/*
#include <glpk.h>
#include <stdlib.h>
#include "mcmf.h"
*/
import "C"
import (
	"strconv"
	"strings"
)

type CustomArc struct {
	Src  int
	Dst  int
	Flow float64
}

type MCMF struct {
	G *C.glp_graph
	P *C.glp_prob
}

func (this *MCMF) New() {
	this.G = C.createGraph();
	this.P = C.createProblem();
}

func (this *MCMF) AddVertices(verticesNum int) {
	C.addVertices(this.G, C.int(verticesNum))
}

func (this *MCMF) SetVertexRhs(idx int, rhs float64) {
	C.setVertexRhs(this.G, C.int(idx), C.double(rhs))
}

func (this *MCMF) AddArc(srcIdx int, dstIdx int, low float64, cap float64, cost float64) {
	C.addArc(this.G, C.int(srcIdx), C.int(dstIdx), C.double(low), C.double(cap), C.double(cost))
}

func (this *MCMF) maxFlowGraph2LinearProblem(srcIdx int, sinkIdx int) {
	if this.G == nil {
		this.P = nil
		return
	}
	C.maxFlowGraph2LinearProblem(this.G, this.P, C.int(srcIdx), C.int(sinkIdx))
}

func (this *MCMF) minCostGraph2LinearProblem() {
	if this.G == nil {
		this.P = nil
		return
	}

	C.minCostGraph2LinearProblem(this.G, this.P)
	//test
	//C.glp_write_lp(this.P, nil, "minCost.lp");
	//C.glp_write_mincost(this.G, C.int(GLP_V_RHS), C.int(GLP_A_LOW), C.int(GLP_A_CAP), C.int(GLP_A_COST), "mincost.dimacs");
}

func (this *MCMF) Solve(srcIdx int, sinkIdx int) []*CustomArc {

	ok := C.solve(this.G, this.P, C.int(srcIdx), C.int(sinkIdx))

	if (bool(ok)) {
		return this.parse2Arcs()
	}
	return nil
}

func (this *MCMF) parse2Arcs() []*CustomArc{
	size := int(C.glp_get_num_cols(this.P))

	arr := make([]*CustomArc, 0, size)

	for i := 0; i < size; i++ {
		name := C.GoString(C.glp_get_col_name(this.P, C.int(i+1)))
		//fmt.Printf("%v, %v, %v, %v\n", name, name[0], string(name[0]), name[2:len(name) -1])

		arc := this.colName2CustomArc(name)
		if arc == nil{
			continue
		}
		arc.Flow = float64(C.glp_get_col_prim(this.P, C.int(i+1)))
		arr = append(arr, arc)
	}
	return arr
}

func (this *MCMF) colName2CustomArc(name string) *CustomArc {
	if name == ""{
		return nil
	}

	if name[0:1]  != "x" || name[1:2] != "[" || name[len(name) -1:] != "]"{
		return nil
	}

	name1 := name[2:len(name) -1]
	arr := strings.Split(name1, ",")
	src, _:= strconv.Atoi(arr[0])
	dst, _:= strconv.Atoi(arr[1])
	arc := CustomArc{Src:src, Dst:dst, Flow:0}
	return &arc
}

