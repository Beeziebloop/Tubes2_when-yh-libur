package backend

import "time"

//nyimpen hasil dari proses traversal
type TraversalRes struct{
	MatchedNodes []*Node //nodes yang ngematch
	VisitedNodes []*Node //nodes visited
	VisCount int //total visited
	ElapsedTime time.Duration
	TraversalLog []StepLog //log detail setiap step traversal
	MaxDepthVisited int //depth terdalam selama transversal saat ini
	FullMaxDepth int //depth maksimum dari seluruh tree
}

//Ngelog sebuah step dalam proses transversal
type StepLog struct{
	Step int
	N_Tag string
	N_ID string
	N_Classes []string
	Depth int
	Is_matched bool
	Algorithm string //bfs atau dfs
}

//nyari maximum depth dari seluruh DOM tree
func maxFullTreeDepth(root *Node) int{
	if root == nil{
		return 0
	}
	maxDepth := 0
	queue := []*Node{root}
	for len(queue) > 0{
		node := queue[0]
		queue = queue[1:]
		if d := node.Depth(); d > maxDepth{
			maxDepth = d
		}
		queue = append(queue, node.Children...)
	}
	return maxDepth
}

//dengan menggunakan referensi dari mediumnya Timothy Britt dan StackOverflow, algoritma ini menggunakan sebuah queue dengan konsep FIFO
func BFS(root *Node, selector string, matchF MatchFunc, topN int, treeMaxDepth int) TraversalRes{
	start := time.Now()
	result := TraversalRes{
		MatchedNodes: []*Node{},
		VisitedNodes: []*Node{},
		TraversalLog: []StepLog{},
		FullMaxDepth: treeMaxDepth,
	}

	if root == nil{
		return result
	}

	queue := []*Node{root}
	visited := make(map[*Node]bool)
	step := 0
	maxDepthVisited := 0

	for len(queue) > 0{
		currentNode := queue[0]
		queue = queue[1:]
		if visited[currentNode]{
			continue
		}
		//mark as visited
		visited[currentNode] = true
		result.VisitedNodes = append(result.VisitedNodes, currentNode)
		step++
		//track visited max depth
		depth := currentNode.Depth()
		if depth > maxDepthVisited{
			maxDepthVisited = depth
		}
		//cek apakah nodenya ngematch sama selector
		matched := matchF(currentNode, selector)
		//catat step kedalam transversal log
		result.TraversalLog = append(result.TraversalLog, StepLog{
			Step: step,
			N_Tag: currentNode.Tag,
			N_ID: currentNode.ID,
			N_Classes: currentNode.Classes,
			Depth: depth,
			Is_matched: matched,
			Algorithm: "bfs",
		})
		if matched{
			result.MatchedNodes = append(result.MatchedNodes, currentNode)
			//berhenti pas jumlah nodes terpenuhi
			if topN != -1 && len(result.MatchedNodes) >= topN{
				break
			}
		}
		//children di enqueue, untuk bfs tidak ada reverse 
		for _, child := range currentNode.Children{
			if !visited[child]{
				queue = append(queue, child)
			}
		}
	}
	result.VisCount = step
	result.ElapsedTime = time.Since(start)
	result.MaxDepthVisited = maxDepthVisited
	return result
}

//dengan menggunakan referensi dari mediumnya Timothy Britt dan StackOverflow, algoritma ini menggunakan sebuah stack LIFO
func DFS(root *Node, selector string, matchF MatchFunc, topN int, treeMaxDepth int) TraversalRes{
	start := time.Now()
	result := TraversalRes{
		MatchedNodes: []*Node{},
		VisitedNodes: []*Node{},
		TraversalLog: []StepLog{},
		FullMaxDepth: treeMaxDepth,
	}

	if root == nil{
		return result
	}

	stack := []*Node{root}
	visited := make(map[*Node]bool)
	step := 0
	maxDepthVisited := 0

	for len(stack) > 0{
		n := len(stack) -1
		currentNode := stack[n]
		stack = stack[:n]
		if visited[currentNode]{
			continue
		}
		//mark as visited
		visited[currentNode] = true
		result.VisitedNodes = append(result.VisitedNodes, currentNode)
		step++
		//track visited max depth
		depth := currentNode.Depth()
		if depth > maxDepthVisited{
			maxDepthVisited = depth
		}
		//cek apakah nodenya ngematch sama selector
		matched := matchF(currentNode, selector)
		//catat stepnya ke traversal log
		result.TraversalLog = append(result.TraversalLog, StepLog{
			Step: step,
			N_Tag: currentNode.Tag,
			N_ID: currentNode.ID,
			N_Classes: currentNode.Classes,
			Depth: depth,
			Is_matched: matched,
			Algorithm: "dfs",
		})
		if matched{
			result.MatchedNodes = append(result.MatchedNodes, currentNode)
			if topN != -1 && len(result.MatchedNodes) >= topN{
				break
			}
		}
		//push children secara reverse biar leftmost di proses terlebih dahulu
		for i := len(currentNode.Children) -1; i >= 0; i--{
			if !visited[currentNode.Children[i]]{
				stack = append(stack, currentNode.Children[i])
			}
		}
	}
	result.VisCount = step
	result.ElapsedTime = time.Since(start)
	result.MaxDepthVisited = maxDepthVisited
	return result
}

type AlgorithmType string
const(
	AlgoBFS AlgorithmType = "bfs"
	AlgoDFS AlgorithmType = "dfs"
)

//main entry point yang dipanggil dari handler, ngekomputasi dulu full max depth tree, baru mulai bfs/dfs
func Search(root *Node, selector string, algoType AlgorithmType, matchF MatchFunc, topN int) TraversalRes{
	maxDepth := maxFullTreeDepth(root)
	switch algoType{
	case AlgoBFS:
		return BFS(root, selector, matchF, topN, maxDepth)
	case AlgoDFS:
		return DFS(root, selector, matchF, topN, maxDepth)
	default:
		return BFS(root, selector, matchF, topN, maxDepth) //defaultnya ke bfs
	}
}