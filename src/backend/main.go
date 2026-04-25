package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// buat nampung data dom tree
type TreeNode struct {
	Tag       string      `json:"tag"`
	ID        string      `json:"id"`
	Classes   []string    `json:"classes"`
	NodeIndex int         `json:"nodeIndex"`
	Children  []*TreeNode `json:"children"`
}

// buat nyatet log tiap langkah
type TraversalLogEntry struct {
	Step      int      `json:"step"`
	NodeTag   string   `json:"nodeTag"`
	NodeId    string   `json:"nodeId"`
	NodeClass []string `json:"nodeClass"`
	Depth     int      `json:"depth"`
	Matched   bool     `json:"matched"`
	NodeIndex int      `json:"nodeIndex"`
}

// SearchRequest dari frontend
type SearchRequest struct {
	URL       string `json:"url"`
	HTML      string `json:"html"`
	Algorithm string `json:"algorithm"`
	Selector  string `json:"selector"`
	TopN      int    `json:"topN"`
}

// SearchResponse buat frontend
type SearchResponse struct {
	Success      bool                `json:"success"`
	Message      string              `json:"message,omitempty"`
	VisitCount   int                 `json:"visitCount"`
	MaxDepth     int                 `json:"maxDepth"`
	ElapsedTime  string              `json:"elapsedTime"`
	MatchedCount int                 `json:"matchedCount"`
	TraversalLog []TraversalLogEntry `json:"traversalLog"`
	Tree         *TreeNode           `json:"tree"`
}

// BuildRequest buat endpoint /api/build
type BuildRequest struct {
	URL  string `json:"url"`
	HTML string `json:"html"`
}

type BuildResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Tree    *TreeNode  `json:"tree"`
}

// global state buat lca
var (
	currentRoot *Node
	currentIndexMap map[*Node]int
	currentNodeByIndex map[int]*Node
	currentLCATable *LCATable
	stateMutex sync.Mutex
)

// buat kasih nomer urut ke tiap node
func assignIndices(root *Node) map[*Node]int {
	indexMap := make(map[*Node]int)
	idx := 0
	var traverse func(n *Node)
	traverse = func(n *Node) {
		if n == nil {
			return
		}
		indexMap[n] = idx
		idx++
		for _, child := range n.Children {
			traverse(child)
		}
	}
	traverse(root)
	return indexMap
}

// ngubah Node jadi TreeNode
func serializeTree(node *Node, indexMap map[*Node]int) *TreeNode {
	if node == nil {
		return nil
	}
	treeNode := &TreeNode{
		Tag: node.Tag,
		ID: node.ID,
		Classes: node.Classes,
		NodeIndex: indexMap[node],
		Children: []*TreeNode{},
	}
	for _, child := range node.Children {
		treeNode.Children = append(treeNode.Children, serializeTree(child, indexMap))
	}
	return treeNode
}

// ngubah StepLog jadi TraversalLogEntry
func convertStepLog(logs []StepLog) []TraversalLogEntry {
	entries := make([]TraversalLogEntry, len(logs))
	for i, l := range logs {
		entries[i] = TraversalLogEntry{
			Step: l.Step,
			NodeTag: l.N_Tag,
			NodeId: l.N_ID,
			NodeClass: l.N_Classes,
			Depth: l.Depth,
			Matched: l.Is_matched,
			NodeIndex: l.NodeIndex,
		}
	}
	return entries
}

func handleBuild(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(BuildResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req BuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(BuildResponse{Success: false, Message: "Invalid JSON"})
		return
	}
	if req.URL == "" && req.HTML == "" {
		json.NewEncoder(w).Encode(BuildResponse{Success: false, Message: "URL or HTML required"})
		return
	}

	var input string
	if req.URL != "" {
		input = req.URL
	} else {
		input = req.HTML
	}

	root, err := LoadHTML(input)
	if err != nil {
		json.NewEncoder(w).Encode(BuildResponse{Success: false, Message: err.Error()})
		return
	}
	if err := ValidateHTML(root); err != nil {
		json.NewEncoder(w).Encode(BuildResponse{Success: false, Message: "HTML structure error: " + err.Error()})
		return
	}

	indexMap := assignIndices(root)
	updateGlobalState(root, indexMap) // simpan state untuk LCA nanti

	tree := serializeTree(root, indexMap)
	json.NewEncoder(w).Encode(BuildResponse{Success: true, Tree: tree})
}

// dipanggil setelah parsing html dan assign indices
func updateGlobalState(root *Node, indexMap map[*Node]int) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	currentRoot = root
	currentIndexMap = indexMap
	// mapping balik
	nodeByIndex := make(map[int]*Node)
	for node, idx := range indexMap {
		nodeByIndex[idx] = node
	}
	currentNodeByIndex = nodeByIndex
	// bangun LCATable
	currentLCATable = buildLCATable(root)
}

// nanganin request lca
func handleLCA(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Method not allowed"})
		return
	}
	var req struct {
		NodeIndexA int `json:"nodeIndexA"`
		NodeIndexB int `json:"nodeIndexB"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Invalid JSON"})
		return
	}

	stateMutex.Lock()
	defer stateMutex.Unlock()
	if currentLCATable == nil || currentNodeByIndex == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Belum ada pohon DOM. Lakukan pencarian terlebih dahulu."})
		return
	}

	nodeA := currentNodeByIndex[req.NodeIndexA]
	nodeB := currentNodeByIndex[req.NodeIndexB]
	if nodeA == nil || nodeB == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Node tidak valid"})
		return
	}

	lcaNode := currentLCATable.lca(nodeA, nodeB)
	if lcaNode == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "LCA tidak ditemukan"})
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"tag":       lcaNode.Tag,
		"id":        lcaNode.ID,
		"classes":   lcaNode.Classes,
		"depth":     currentLCATable.depth[lcaNode],
		"nodeIndex": currentIndexMap[lcaNode],
	}
	json.NewEncoder(w).Encode(response)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// set header json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(SearchResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	// validasi input, harus ada url atau html
	if req.URL == "" && req.HTML == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResponse{Success: false, Message: "Either url or html must be provided"})
		return
	}

	// load html dari input
	var input string
	if req.URL != "" {
		input = req.URL
	} else {
		input = req.HTML
	}
	root, err := LoadHTML(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResponse{Success: false, Message: "Failed to load HTML: " + err.Error()})
		return
	}

	// validasi struktur html
	if err := ValidateHTML(root); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResponse{Success: false, Message: "HTML structure error: " + err.Error()})
		return
	}

	// bikin index map
	indexMap := assignIndices(root)

	// simpen state lca
	updateGlobalState(root, indexMap)

	// bikin match function
	matchFunc := MakeMatchFunc(req.Selector)

	// pilih mau pake cara dfs atau bfs
	var algoType AlgorithmType
	switch req.Algorithm {
	case "DFS":
		algoType = AlgoDFS
	default:
		algoType = AlgoBFS
	}

	// lakuin pencarian
	result := Search(root, req.Selector, algoType, matchFunc, req.TopN, indexMap)

	// bikin struktur pohon
	tree := serializeTree(root, indexMap)

	// konversi log
	logEntries := convertStepLog(result.TraversalLog)

	// hitung matched count
	matchedCount := len(result.MatchedNodes)

	// format elapsed time
	elapsedMs := float64(result.ElapsedTime.Nanoseconds()) / 1e6
	elapsedStr := fmt.Sprintf("%.3f ms", elapsedMs)

	response := SearchResponse{
		Success:      true,
		VisitCount:   result.VisCount,
		MaxDepth:     result.MaxDepthVisited,
		ElapsedTime:  elapsedStr,
		MatchedCount: matchedCount,
		TraversalLog: logEntries,
		Tree:         tree,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	frontendDir := "../frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		log.Fatalf("Frontend directory not found at %s", frontendDir)
	}

	// handler file statis
	fs := http.FileServer(http.Dir(frontendDir))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))

	// api endpoint
	http.HandleFunc("/api/search", handleSearch)
	http.HandleFunc("/api/build", handleBuild)
	http.HandleFunc("/api/lca", handleLCA)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}