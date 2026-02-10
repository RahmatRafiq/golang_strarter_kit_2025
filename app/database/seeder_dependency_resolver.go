package database

import (
	"fmt"
)

// DependencyGraph represents the dependency relationships between seeders
type DependencyGraph struct {
	nodes map[string]*Seeder
	edges map[string][]string
}

// NewDependencyGraph creates a new dependency graph from a list of seeders
func NewDependencyGraph(seeders []Seeder) *DependencyGraph {
	dg := &DependencyGraph{
		nodes: make(map[string]*Seeder),
		edges: make(map[string][]string),
	}

	for i := range seeders {
		s := &seeders[i]
		dg.nodes[s.Name] = s
		dg.edges[s.Name] = s.DependsOn
	}

	return dg
}

// TopologicalSort performs a topological sort on the dependency graph
// Returns the seeders in the correct order to satisfy all dependencies
func (dg *DependencyGraph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	result := []string{}

	var visit func(string, []string) error
	visit = func(node string, path []string) error {
		// Check for circular dependency
		if recStack[node] {
			cyclePath := append(path, node)
			return fmt.Errorf("circular dependency detected: %v", cyclePath)
		}

		// Already processed
		if visited[node] {
			return nil
		}

		// Mark as being processed
		recStack[node] = true
		newPath := append(path, node)

		// Process dependencies first
		for _, dep := range dg.edges[node] {
			// Check if dependency exists
			if _, exists := dg.nodes[dep]; !exists {
				return fmt.Errorf("dependency '%s' not found for seeder '%s'", dep, node)
			}

			if err := visit(dep, newPath); err != nil {
				return err
			}
		}

		// Done processing this node
		recStack[node] = false
		visited[node] = true

		// Add to result (dependencies come first)
		result = append([]string{node}, result...)

		return nil
	}

	// Visit all nodes
	for node := range dg.nodes {
		if err := visit(node, []string{}); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ValidateDependencies checks if all dependencies are valid
func (dg *DependencyGraph) ValidateDependencies() error {
	for seederName, deps := range dg.edges {
		for _, dep := range deps {
			if _, exists := dg.nodes[dep]; !exists {
				return fmt.Errorf("seeder '%s' depends on '%s' which does not exist", seederName, dep)
			}
		}
	}
	return nil
}

// GetDependencyChain returns the full dependency chain for a seeder
func (dg *DependencyGraph) GetDependencyChain(seederName string) ([]string, error) {
	if _, exists := dg.nodes[seederName]; !exists {
		return nil, fmt.Errorf("seeder '%s' not found", seederName)
	}

	visited := make(map[string]bool)
	chain := []string{}

	var collect func(string) error
	collect = func(node string) error {
		if visited[node] {
			return nil
		}

		visited[node] = true

		// Collect dependencies first
		for _, dep := range dg.edges[node] {
			if err := collect(dep); err != nil {
				return err
			}
		}

		chain = append(chain, node)
		return nil
	}

	if err := collect(seederName); err != nil {
		return nil, err
	}

	return chain, nil
}
