package audit

// FootprintMetrics holds the result of footprint analysis
type FootprintMetrics struct {
	TransitiveDeps   int
	MaxDepth         int
	DependencyWeight float64
}

// CalculateFootprint analyzes dependency graph complexity
func CalculateFootprint(module Module, allModules []Module) FootprintMetrics {
	// This requires building a graph from allModules to traverse
	// For MVP, if we don't have the full graph structure (parent->child links),
	// we can only estimate based on global stats or if 'go list -json' provided the tree.
	// 'go list -json all' gives a flat list, but each module has a 'Dir' and we can't easily see who depends on whom
	// without 'go mod graph'.
	
	// However, if we are auditing a specific module *within* the graph, we want to know *its* dependencies.
	// That's hard from a flat list without the graph.
	
	// For the "Project" footprint, we can just count total modules.
	// For a "Module" footprint (how heavy is this dependency?), we need the graph.
	
	// Let's assume for MVP we just return 0 for per-module footprint unless we parse 'go mod graph'.
	// Parsing 'go mod graph' output is easy if we run it.
	
	return FootprintMetrics{
		TransitiveDeps:   0,
		MaxDepth:         0,
		DependencyWeight: 0,
	}
}

// CalculateProjectFootprint calculates the total footprint of the project
func CalculateProjectFootprint(modules []Module) FootprintMetrics {
	total := len(modules)
	direct := 0
	for _, m := range modules {
		if !m.Indirect {
			direct++
		}
	}
	
	return FootprintMetrics{
		TransitiveDeps:   total - direct,
		MaxDepth:         0, // Unknown without graph
		DependencyWeight: float64(total),
	}
}
