package main

import (
	_ "embed"
	"eventflowdb/projections"
	"log"
)

func sandbox() {
	compiler, err := projections.NewCompiler()
	check(err)

	code, err := compiler.Compile(`
		project({
			ProductAdded: (state, event) => {
				state.name = event.name;
			}
		});
	`)
	check(err)

	log.Println(code)
}
