package main
//algoritma lca menggunakan binary lifting ini mengambil referensi dari https://cp-algorithms.com/graph/lca_binary_lifting.html serta mengambil inspirasi dari https://www.youtube.com/watch?v=oib-XsjFa-M
//berikut adalah tahapannya preprocessing untuk binary lifting:
//1. dfs traversal dulu untuk mencatat waktu masuk (tin) dan waktu keluar (tout) tiap node
//2. buat tabel ancestor dimana nantinya ancestor[v][k] adalah ancestor ke 2^k dari v
//untuk tiap query lca(u, v):
//1. cek dulu apakah salah satu dari u atau v adalah ancestor dari yang lain, kalau iya itu lcanya
//2. kalau bukan, naikin ke atas sebesar 2^k langkah, berhenti tepat sebelum ancestor u dan v jadi saudara
//3. lca = ancestor[u][0] (direct parent u setelah loop)

const maxLog = 64 //ini sekedar buat safety cap aja

//LCATable ini untuk nampung semua info preprocessed untuk query lca nanti
type LCATable struct{
	tIn map[*Node]int //waktu masuk tiap node
	tOut map[*Node]int //waktu keluar tiap node
	depth map[*Node]int //depth tiap node
	ancestor map[*Node][]*Node //ancestor[v][k] itu ancestor ke 2^k dari v
	nodeList []*Node //nyimpan semua node dalam urutan dfs untuk preprocessing
	timer int //global dfs timestamp
	log int //the computed log2(n)
}

//buat ngitung ceil(log2(n)) secara gabole pake library math, minimal 1
func computeLog(n int) int{
	log := 1
	for (1 << log) <= n{
		log++
		if log >= maxLog{
			break
		}
	}
	return log
}

func (t *LCATable) dfs(v *Node, parent *Node){
	t.timer++;
	t.tIn[v] = t.timer
	t.nodeList = append(t.nodeList, v)
	//inisialisasi slice ancestor[v][0] dengan parent langsungnya
	anc := make([]*Node, 1)
	if parent != nil{
		anc[0] = parent
		t.depth[v] = t.depth[parent] + 1
	}else{
		anc[0] = nil //root gada parent
		t.depth[v] = 0
	}
	t.ancestor[v] = anc

	for _, child := range v.Children{
		t.dfs(child, v)
	}

	t.timer++
	t.tOut[v] = t.timer
}

//ngepreprocess dom tree dari root dan mengembalikan LCATable, harus dipanggil sekali setelah parseHTML
func buildLCATable(root *Node) *LCATable{
	table := &LCATable{
		tIn: make(map[*Node]int),
		tOut: make(map[*Node]int),
		depth: make(map[*Node]int),
		ancestor: make(map[*Node][]*Node),
		nodeList: []*Node{},
		timer: 0,
		log: 1, //akan diupdate setelah dfs 
	}

	//dfs dulu buat isi tIn, tOut, depth, dan ancestor[v][0]
	table.dfs(root, nil)

	//lalu hitung log berdasarkan jumlah node
	n := len(table.nodeList)
	table.log = computeLog(n)

	//expand semua slice ancestor hingga panjang log full, terus isi ancestor[v][k] untuk k = 1...log-1
	//ancestor[v][k] = ancestor[ancestor[v][k-1]][k-1]
	for _, node := range table.nodeList{
		anc := table.ancestor[node]
		//kembangin slicenya jadi panjang log (karena sebelumnya diinisialisasi panjangnya 1 dari dfs)
		for len(anc) < table.log{
			anc = append(anc, nil)
		}
		table.ancestor[node] = anc
	}
	for k := 1; k < table.log; k++{
		for _, node := range table.nodeList{
			anc := table.ancestor[node]
			if anc[k-1] != nil{
				parAnc := table.ancestor[anc[k-1]]
				if k - 1 < len(parAnc){
					anc[k] = parAnc[k-1]
				}
			}
			table.ancestor[node] = anc
		}
	}
	
	return table
}

//ancestor validator
func (t *LCATable) isAnc(u, v *Node) bool{
	return t.tIn[u] <= t.tIn[v] && t.tOut[u] >= t.tOut[v]
}

//IsAncestor ini mengekspos fungsi isAnc secara publik untuk API handler
func (t *LCATable) IsAncestor(u, v *Node) bool{
	return t.isAnc(u, v)
}

func (t *LCATable) lca(u, v *Node) *Node{
	if u == nil || v == nil{
		return nil
	}
	//cek dulu apakah kedua node exist dalam tabel
	if _, ok := t.tIn[u]; !ok{
		return nil
	}
	if _, ok := t.tIn[v]; !ok{
		return nil
	}

	//cek apakah salah satu node adalah ancestor dari yang lain
	if t.isAnc(u, v){
		return u
	}
	if t.isAnc(v, u){
		return v
	}

	//naikkan u ke atas sebesar 2^k langkah sampai tepat sebelum ancestor u dan v jadi saudara (cari node tertinggi yang bukan ancestor si v)
	for k := t.log - 1; k >= 0; k--{
		anc := t.ancestor[u]
		if k < len(anc) && anc[k] != nil && !t.isAnc(anc[k], v){
			u = anc[k]
		}
	}

	//lca adalah direct parent dari u setelah loop
	return t.ancestor[u][0]
}

func (t *LCATable) nodeDepth(n *Node) int{
	return t.depth[n]
}

//ngitung jumlah edges antara u dan v
func (t *LCATable) nodeDistance(u, v *Node) int{
	lca := t.lca(u, v)
	if lca == nil{
		return -1
	}
	return t.depth[u] + t.depth[v] - (2*t.depth[lca])
}

//ngembaliin k-th ancestor dari node n (1 untuk parent, 2 untuk grandparent, etc). ngembaliin nil kalau k melebihi dari depth node
func (t *LCATable) kthAnc(n *Node, k int) *Node{
	if n == nil || k < 0{
		return nil
	}
	if t.depth[n] < k{
		return nil
	}
	for j := t.log - 1; j >= 0; j--{
		if k >= (1 << j){
			anc := t.ancestor[n]
			if j < len(anc) && anc[j] == nil{
				return nil
			}
			n = anc[j]
			k -= 1 << j
		}
	}
	return n
}