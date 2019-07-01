package glpk

// #cgo LDFLAGS: -lglpk
// #include <glpk.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

type graph struct {
	g *C.glp_graph
}

type Graph struct {
	g *graph
}

type vertex struct {
	v *C.glp_vertex
}

type Vertex struct {
	v *vertex
}

type arc struct {
	a *C.glp_arc
}

type Arc struct {
	a *arc
}

/**
create new graph
*/
func NewGraph(v_size int, a_size int) *Graph {
	g := &graph{C.glp_create_graph(v_size, a_size)}
	return &Graph{g}
}

// set graph name
func (g *Graph) SetGraphName(name string) {
	if g.g.g == nil {
		panic("Graph method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_graph_name(g.g.g, s)
}

// add vertices
func (g *Graph) AddVertices(nadd int) int {
	if g.g.g == nil {
		panic("Graph method called on a deleted problem")
	}

	return int(C.glp_add_vertices(g.g.g, C.int(nadd)))
}

// set vertex name
func (g *Graph) SetVertexName(i int, name string) {
	if g.g.g == nil {
		panic("Graph method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	C.glp_set_vertex_name(g.g.g, C.int(i), s)

}

// add arc
func (g *Graph) AddArc(i int, j int) *Arc {
	if g.g.g == nil {
		panic("Graph method called on a deleted problem")
	}
	a := &arc{C.glp_add_arc(g.g.g, C.int(i), C.int(j))}
	return &Arc{a}
}

/* delete arc from graph */
func (g *Graph) DelArc(a *Arc){
	if g.g.g == nil {
		panic("Graph method called on a deleted problem")
	}
	if  a.a.a == nil {
		panic("Arc method called on a deleted problem")
	}

	C.glp_del_arc(g.g.g, a.a.a)
}

/* delete vertices from graph */
func (g *Graph) DelVertices(ndel int, num []int){
	if g.g.g == nil {
		panic("Graph method called on a deleted problem")
	}
	// todo
}


/* delete graph */
func (g *Graph) DeleteGraph(){
	if g.g.g == nil {
		panic("Graph method called on a deleted problem")
	}
	C.glp_delete_graph(g.g.g)
}

