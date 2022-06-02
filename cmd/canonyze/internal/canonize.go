package internal

import (
	"fmt"
	"io"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

func Canonize() error {
	err := canonize(os.Stdin, os.Stdout)
	if err != nil {
		return fmt.Errorf("canonizing stdin to stdout: %w", err)
	}
	return nil
}

func canonize(reader io.Reader, writer io.Writer) error {
	docs, err := readDocs(reader)
	if err != nil {
		return fmt.Errorf("reading documents from yaml: %w", err)
	}

	canonizeDocs(docs)

	err = writeDocs(writer, docs)
	if err != nil {
		return fmt.Errorf("writing documents to yaml: %w", err)
	}

	return nil
}

func canonizeDocs(docs []*yaml.Node) error {
	sortDocsByKindThenMetadataName(docs)
	for _, doc := range docs {
		for _, content := range doc.Content {
			sortNodeContent(content)
		}
	}
	return nil
}

func sortDocsByKindThenMetadataName(docs []*yaml.Node) {
	sort.Slice(docs, func(i, j int) bool {
		doc1 := docs[i].Content[0]
		doc2 := docs[j].Content[0]
		kind1 := getMappingScalarByName(doc1, "kind")
		kind2 := getMappingScalarByName(doc2, "kind")
		if kind1 == kind2 {
			name1 := getMetadataName(doc1)
			name2 := getMetadataName(doc2)
			return name1 < name2
		}
		return kind1 < kind2
	})
}

type entry struct {
	key, value *yaml.Node
}

func sortMapping(nodes []*yaml.Node) {
	// Convert nodes to entries for sorting
	entries := make([]entry, len(nodes)/2)
	for i := 0; i < len(nodes)/2; i++ {
		entries[i].key = nodes[i*2]
		entries[i].value = nodes[i*2+1]
	}

	// Sort entries
	sort.Slice(entries, func(i, j int) bool {
		name1 := entries[i].key.Value
		name2 := entries[j].key.Value
		return name1 < name2
	})

	// Convert entries back to nodes
	for i := 0; i < len(entries); i++ {
		nodes[i*2] = entries[i].key
		nodes[i*2+1] = entries[i].value
	}

	// Sort content of nodes recursively
	for _, node := range nodes {
		sortNodeContent(node)
	}
}

func sortNodeContent(node *yaml.Node) {
	switch node.Kind {
	case yaml.SequenceNode:
		for _, content := range node.Content {
			sortNodeContent(content)
		}
	case yaml.MappingNode:
		sortMapping(node.Content)
	}
}

func getMetadataName(doc *yaml.Node) string {
	metadata := getMappingNodeByName(doc, "metadata")
	if metadata == nil {
		return ""
	}

	return getMappingScalarByName(metadata, "name")
}

func getMappingScalarByName(node *yaml.Node, name string) string {
	index := findMappingItemByName(node, name)
	if index == -1 {
		return ""
	}
	return node.Content[index*2+1].Value
}

func getMappingNodeByName(node *yaml.Node, name string) *yaml.Node {
	index := findMappingItemByName(node, name)
	if index == -1 {
		return nil
	}
	return node.Content[index*2+1]
}

// findMappingItemByName returns index of  (numbered 0 for first key/value pair,
// 1 for second, and so on) or -1 if not found.
func findMappingItemByName(node *yaml.Node, name string) int {
	for i := 0; i < len(node.Content)/2; i++ {
		if node.Content[i*2].Value == name {
			return i
		}
	}
	return -1
}

func readDocs(reader io.Reader) ([]*yaml.Node, error) {
	var nodes []*yaml.Node
	decoder := yaml.NewDecoder(reader)
	for {
		var node yaml.Node
		if err := decoder.Decode(&node); err != nil {
			// Reached end of documents to decode?
			if err == io.EOF {
				return nodes, nil
			}
			return nil, fmt.Errorf("decoding document node: %w", err)
		}
		nodes = append(nodes, &node)
	}
}

func writeDocs(writer io.Writer, nodes []*yaml.Node) error {
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	for _, node := range nodes {
		if err := encoder.Encode(node); err != nil {
			return fmt.Errorf("encoding document node: %w", err)
		}
	}
	return nil
}
