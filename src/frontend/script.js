// ========== GLOBAL STATE ==========
let currentMode = 'css'; // 'css' or 'lca'
let selectedNodeA = null;
let selectedNodeB = null;
let lcaActive = true;

// ========== DOM ELEMENTS ==========
const buildTreeBtn = document.getElementById('buildTreeBtn');
const searchBtn = document.getElementById('searchBtn');
const findLCABtn = document.getElementById('findLCABtn');
const resetLCASelectionBtn = document.getElementById('resetLCASelectionBtn');
const cssPanel = document.getElementById('cssPanel');
const lcaPanel = document.getElementById('lcaPanel');
const selectedASpan = document.getElementById('selectedA');
const selectedBSpan = document.getElementById('selectedB');
const lcaResultArea = document.getElementById('lcaResultArea');

// ========== UTILITY FUNCTIONS ==========

// buat ganti tampilan input mau pake url atau raw html
function switchInputSource(value) {
    const urlContainer = document.getElementById('input-url-container');
    const rawContainer = document.getElementById('input-raw-container');

    if (value === 'url') {
        urlContainer.classList.remove('hidden');
        rawContainer.classList.add('hidden');
    } else {
        urlContainer.classList.add('hidden');
        rawContainer.classList.remove('hidden');
    }
}

// buat ganti mode css selector atau lca
function setAppMode(mode) {
    currentMode = mode;
    const btnCss = document.getElementById('btn-mode-css');
    const btnLca = document.getElementById('btn-mode-lca');

    if (mode === 'css') {
        btnCss.classList.add('active');
        btnLca.classList.remove('active');
        cssPanel.style.display = 'block';
        lcaPanel.style.display = 'none';
        buildTreeBtn.style.display = 'none';
    } else {
        btnCss.classList.remove('active');
        btnLca.classList.add('active');
        cssPanel.style.display = 'none';
        lcaPanel.style.display = 'block';
        buildTreeBtn.style.display = 'block';
        
        resetLCA();
        lcaActive = false;
        disableNodeSelection();
    }
}

// buat reset semua pilihan node lca
function resetLCA() {
    selectedNodeA = null;
    selectedNodeB = null;

    if (selectedASpan) selectedASpan.innerText = '-';
    if (selectedBSpan) selectedBSpan.innerText = '-';
    if (findLCABtn) findLCABtn.disabled = true;
    if (lcaResultArea) lcaResultArea.innerHTML = '';

    // hapus highlight warna di domtree
    document.querySelectorAll('.node-element').forEach(el => {
        el.classList.remove('selected-a', 'selected-b', 'lca-node');
    });

    if (currentMode === 'lca') {
        lcaActive = true;
        enableNodeSelection();
    }
}

// buat update warna node yang lagi dipilih
function updateNodeHighlight() {
    document.querySelectorAll('.node-element').forEach(el => {
        el.classList.remove('selected-a', 'selected-b');
    });

    if (selectedNodeA !== null) {
        const elA = document.querySelector(`.node-element[data-index="${selectedNodeA}"]`);
        if (elA) elA.classList.add('selected-a');
    }

    if (selectedNodeB !== null) {
        const elB = document.querySelector(`.node-element[data-index="${selectedNodeB}"]`);
        if (elB) elB.classList.add('selected-b');
    }
}

// pasang event listener klik ke tiap elemen node di tree
function enableNodeSelection() {
    if (!lcaActive) return;
    const nodes = document.querySelectorAll('.node-element');
    nodes.forEach(node => {
        node.removeEventListener('click', nodeClickHandler);
        node.addEventListener('click', nodeClickHandler);
    });
}

// cabut event listener klik dari elemen node
function disableNodeSelection() {
    const nodes = document.querySelectorAll('.node-element');
    nodes.forEach(node => {
        node.removeEventListener('click', nodeClickHandler);
    });
}

// logika pas user klik node di tree pas mode lca aktif
function nodeClickHandler(e) {
    if (currentMode !== 'lca' || !lcaActive) return;
    e.stopPropagation();

    const targetSpan = e.currentTarget;
    const nodeIndex = parseInt(targetSpan.dataset.index);
    if (isNaN(nodeIndex)) return;

    // unselect kalo klik node yang udah kepilih
    if (selectedNodeA === nodeIndex) {
        selectedNodeA = null;
        selectedASpan.innerText = '-';
        findLCABtn.disabled = true;
        updateNodeHighlight();
        return;
    }

    if (selectedNodeB === nodeIndex) {
        selectedNodeB = null;
        selectedBSpan.innerText = '-';
        findLCABtn.disabled = (selectedNodeA === null);
        updateNodeHighlight();
        return;
    }

    // isi slot node
    if (selectedNodeA === null) {
        selectedNodeA = nodeIndex;
        selectedASpan.innerText = `Node ${nodeIndex}`;
        updateNodeHighlight();
    } else if (selectedNodeB === null) {
        if (selectedNodeA !== nodeIndex) {
            selectedNodeB = nodeIndex;
            selectedBSpan.innerText = `Node ${nodeIndex}`;
            updateNodeHighlight();
            findLCABtn.disabled = false;
        } else {
            alert('Node sudah dipilih sebagai A. Pilih node lain untuk B.');
        }
    } else {
        alert('Kedua node sudah terpilih. Gunakan "Reset Pilihan" untuk memilih ulang.');
    }
}

// biar layar auto scroll ke node yang dimaksud
function scrollToNode(nodeElement) {
    if (!nodeElement) return;
    const container = document.querySelector('.canvas-container');
    if (!container) return;

    const containerRect = container.getBoundingClientRect();
    const nodeRect = nodeElement.getBoundingClientRect();
    const scrollTop = container.scrollTop;
    
    const targetTop = scrollTop + nodeRect.top - containerRect.top - (containerRect.height / 2) + (nodeRect.height / 2);
    
    container.scrollTo({ top: targetTop, behavior: 'smooth' });
}

// ========== RENDER FUNCTIONS ==========

// buat ngubah data json jadi domtree
function renderTree(node, depth = 0) {
    if (!node) return '';
    const idx = node.nodeIndex;
    const classNames = node.classes ? node.classes.join(' .') : '';
    const idAttr = node.id ? ` #${node.id}` : '';
    
    let attrsHtml = '';
    if (node.attributes) {
        const attrs = Object.entries(node.attributes);
        if (attrs.length > 0) {
            attrsHtml = ' <span class="node-attrs">[' + attrs.map(([k,v]) => `${k}="${v}"`).join(', ') + ']</span>';
        }
    }

    const tagDisplay = `<span class="node-element" data-index="${idx}">
        <span class="node-index">[${idx}]</span>
        <strong>${node.tag}</strong>${idAttr}${classNames ? ` .${classNames}` : ''}
        ${attrsHtml}
    </span>`;

    let childrenHtml = '';
    if (node.children && node.children.length) {
        childrenHtml = '<ul>' + node.children.map(child => renderTree(child, depth + 1)).join('') + '</ul>';
    }

    return `<li>${tagDisplay}${childrenHtml}</li>`;
}

// buat nampilin log traversal
function renderLog(entries) {
    return entries.map(entry => {
        const badge = entry.matched ? '<span class="matched-badge">✓ MATCH</span>' : '';
        const classStr = entry.nodeClass?.join(' ') || '-';
        let attrsStr = '-';

        if (entry.nodeAttributes && Object.keys(entry.nodeAttributes).length > 0) {
            attrsStr = Object.entries(entry.nodeAttributes).map(([k,v]) => `${k}="${v}"`).join(', ');
        }

        return `<li>
            <strong>Step ${entry.step}</strong> | 
            Tag: ${entry.nodeTag} | 
            ID: ${entry.nodeId || '-'} | 
            Class: ${classStr} | 
            Attributes: ${attrsStr} | 
            Depth: ${entry.depth} | 
            Index: ${entry.nodeIndex}
            ${badge}
        </li>`;
    }).join('');
}

// ========== PATH HIGHLIGHTING ==========

// hapus warna jalur root ke target
function resetPaths() {
    document.querySelectorAll('.node-element').forEach(el => {
        el.classList.remove('path-node');
    });
    document.querySelectorAll('.canvas-container li').forEach(li => {
        li.classList.remove('path-li');
    });
}

// kasih highlight jalur dari node yang match sampe root
function highlightPaths(logs) {
    if (!logs || logs.length === 0) return;
    const matchedIndices = new Set();
    
    logs.forEach(entry => {
        if (entry.matched) matchedIndices.add(entry.nodeIndex);
    });

    if (matchedIndices.size === 0) return;

    // telusurin ke atas lewat parent dom elemen
    matchedIndices.forEach(idx => {
        const targetSpan = document.querySelector(`.node-element[data-index="${idx}"]`);
        if (!targetSpan) return;

        let currentLi = targetSpan.closest('li');
        while (currentLi) {
            currentLi.classList.add('path-li');
            const nodeSpan = currentLi.querySelector('.node-element');
            if (nodeSpan) nodeSpan.classList.add('path-node');
            
            // cari parent <li> terdekat di atasnya
            const parentUl = currentLi.parentElement?.closest('li');
            currentLi = parentUl;
        }
    });
}

// ========== RESET HIGHLIGHTS ==========
// balikin semua elemen ke warna normal (hapus status visited/match)
function resetHighlights() {
    document.querySelectorAll('.node-element').forEach(el => {
        el.classList.remove('visited', 'matched', 'current-step', 'path-node');
    });
    document.querySelectorAll('.canvas-container li').forEach(li => {
        li.classList.remove('visited', 'path-li');
    });
}

// ========== ANIMASI ==========

// buat anismasi
async function animateTraversal(logs) {
    const delay = 300; // Kecepatan animasi (ms)

    for (let i = 0; i < logs.length; i++) {
        const step = logs[i];
        const logEntries = document.querySelectorAll('#logContainer li');
        
        if (logEntries[i]) logEntries[i].classList.add('current-step');

        const targetNode = document.querySelector(`.node-element[data-index="${step.nodeIndex}"]`);
        if (targetNode) {
            const parentLi = targetNode.closest('li');
            targetNode.classList.add('visited');
            if (parentLi) parentLi.classList.add('visited');
            if (step.matched) targetNode.classList.add('matched');
            
            targetNode.classList.add('current-step');
            scrollToNode(targetNode);
        }

        await new Promise(resolve => setTimeout(resolve, delay));
        
        if (targetNode) targetNode.classList.remove('current-step');
        if (logEntries[i]) logEntries[i].classList.remove('current-step');
    }

    highlightPaths(logs);
}

// highlight warna
function applyAllHighlights(logs) {
    for (const step of logs) {
        const targetNode = document.querySelector(`.node-element[data-index="${step.nodeIndex}"]`);
        if (targetNode) {
            const parentLi = targetNode.closest('li');
            targetNode.classList.add('visited');
            if (parentLi) parentLi.classList.add('visited');
            if (step.matched) targetNode.classList.add('matched');
        }
    }
    highlightPaths(logs);
}

// ========== API CALLS ==========

// tombol build tree
buildTreeBtn.addEventListener('click', async () => {
    if (currentMode !== 'lca') return;
    const isUrlMode = document.getElementById('inputSourceSelect').value === 'url';
    let body = {};

    if (isUrlMode) {
        const url = document.getElementById('url').value.trim();
        if (!url) return alert('Masukkan URL!');
        body = { url: url };
    } else {
        const html = document.getElementById('html-raw').value.trim();
        if (!html) return alert('Masukkan kode HTML!');
        body = { html: html };
    }

    document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Memproses...</p>';
    
    try {
        const res = await fetch('http://localhost:8080/api/build', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        const data = await res.json();

        if (!data.success) {
            alert('Gagal build DOM: ' + (data.message || 'Unknown error'));
            return;
        }

        document.getElementById('treeContainer').innerHTML = '<ul>' + renderTree(data.tree) + '</ul>';
        resetLCA();
        lcaActive = true;
        enableNodeSelection();
    } catch (err) {
        console.error(err);
        alert('Gagal menghubungi server: ' + err.message);
    }
});

// tombol search css selector
searchBtn.addEventListener('click', async () => {
    if (currentMode !== 'css') {
        alert('Pilih mode CSS Selector terlebih dahulu.');
        return;
    }

    const algorithm = document.getElementById('algorithm').value;
    const selector = document.getElementById('selector').value.trim();
    let topN = parseInt(document.getElementById('topN').value, 10);
    const isUrlMode = document.getElementById('inputSourceSelect').value === 'url';
    const enableAnimation = document.getElementById('enableAnimation').checked;

    // validasi input user
    if (isNaN(topN)) topN = -1;
    if (topN < -1) {
        alert("Top N tidak boleh kurang dari -1.");
        document.getElementById('topN').value = -1;
        topN = -1;
    }
    if (topN == 0) {
        alert("Tidak melakukan percarian apapun. Pilih angka lain!");
        document.getElementById('topN').value = -1;
        topN = -1;
    }

    if (!selector) {
        alert("CSS Selector tidak boleh kosong!");
        return;
    }

    let body = { algorithm, selector, topN };
    if (isUrlMode) {
        const url = document.getElementById('url').value.trim();
        if (!url) return alert('Masukkan URL!');
        body.url = url;
    } else {
        const html = document.getElementById('html-raw').value.trim();
        if (!html) return alert('Masukkan kode HTML!');
        body.html = html;
    }

    document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Memproses...</p>';
    document.getElementById('logContainer').innerHTML = '<li class="hint">Memuat log...</li>';

    try {
        const res = await fetch('http://localhost:8080/api/search', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        const data = await res.json();

        if (!data.success) throw new Error(data.message || 'Unknown error');

        // update stats
        document.getElementById('visitCount').innerText = data.visitCount ?? '-';
        document.getElementById('maxDepth').innerText = data.maxDepth ?? '-';
        document.getElementById('timeMs').innerText = data.elapsedTime ?? '-';
        document.getElementById('matchedCount').innerText = data.matchedCount ?? 0;

        if (data.tree) {
            document.getElementById('treeContainer').innerHTML = '<ul>' + renderTree(data.tree) + '</ul>';
        } else {
            document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Tidak ada DOM tree</p>';
        }

        // reset tampilan lama sebelum nampilin yang baru
        resetLCA(); 
        resetHighlights();
        resetPaths(); 

        if (data.traversalLog && data.traversalLog.length) {
            document.getElementById('logContainer').innerHTML = renderLog(data.traversalLog);
            if (enableAnimation) {
                await animateTraversal(data.traversalLog);
            } else {
                applyAllHighlights(data.traversalLog);
            }
        } else {
            document.getElementById('logContainer').innerHTML = '<li class="hint">Tidak ada log traversal.</li>';
        }
    } catch (err) {
        console.error(err);
        alert('Gagal: ' + err.message);
        document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Error</p>';
    }
});

// tombol search lca
findLCABtn.addEventListener('click', async () => {
    if (selectedNodeA === null || selectedNodeB === null) {
        alert('Pilih dua node terlebih dahulu.');
        return;
    }

    try {
        const res = await fetch('http://localhost:8080/api/lca', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ nodeIndexA: selectedNodeA, nodeIndexB: selectedNodeB })
        });
        const data = await res.json();

        if (!data.success) {
            lcaResultArea.innerHTML = `<span style="color:red;">Error: ${data.message}</span>`;
            return;
        }

        const classStr = data.classes ? data.classes.join(' ') : '';
        lcaResultArea.innerHTML = `<strong>LCA:</strong> &lt;${data.tag}&gt;${data.id ? '#'+data.id : ''}${classStr ? ' .'+classStr : ''} (Depth ${data.depth}, Index ${data.nodeIndex})`;
        
        // kasih tanda node mana yang LCA
        document.querySelectorAll('.node-element').forEach(el => el.classList.remove('lca-node'));
        const lcaSpan = document.querySelector(`.node-element[data-index="${data.nodeIndex}"]`);
        
        if (lcaSpan) {
            lcaSpan.classList.add('lca-node');
            scrollToNode(lcaSpan);
        }

        lcaActive = false;
        disableNodeSelection();
        findLCABtn.disabled = true;
    } catch (err) {
        console.error(err);
        lcaResultArea.innerHTML = `<span style="color:red;">Gagal menghubungi server: ${err.message}</span>`;
    }
});

// tombol reset pilihan node di lca
resetLCASelectionBtn.addEventListener('click', () => {
    resetLCA();
    updateNodeHighlight();
    findLCABtn.disabled = true;
    lcaResultArea.innerHTML = '';

    if (currentMode === 'lca') {
        lcaActive = true;
        enableNodeSelection();
    }
});