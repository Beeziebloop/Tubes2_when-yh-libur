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

function resetLCA() {
    selectedNodeA = null;
    selectedNodeB = null;
    if (selectedASpan) selectedASpan.innerText = '-';
    if (selectedBSpan) selectedBSpan.innerText = '-';
    if (findLCABtn) findLCABtn.disabled = true;
    if (lcaResultArea) lcaResultArea.innerHTML = '';

    document.querySelectorAll('.node-element').forEach(el => {
        el.classList.remove('selected-a', 'selected-b', 'lca-node');
    });

    if (currentMode === 'lca') {
        lcaActive = true;
        enableNodeSelection();
    }
}

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

function enableNodeSelection() {
    if (!lcaActive) return;
    const nodes = document.querySelectorAll('.node-element');
    nodes.forEach(node => {
        node.removeEventListener('click', nodeClickHandler);
        node.addEventListener('click', nodeClickHandler);
    });
}

function disableNodeSelection() {
    const nodes = document.querySelectorAll('.node-element');
    nodes.forEach(node => {
        node.removeEventListener('click', nodeClickHandler);
    });
}

function nodeClickHandler(e) {
    if (currentMode !== 'lca' || !lcaActive) return;
    e.stopPropagation();
    const targetSpan = e.currentTarget;
    const nodeIndex = parseInt(targetSpan.dataset.index);
    if (isNaN(nodeIndex)) return;

    // toggle
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

function resetHighlights() {
    document.querySelectorAll('.node-element').forEach(el => {
        el.classList.remove('visited', 'matched', 'current-step');
    });
    document.querySelectorAll('.canvas-container li').forEach(li => {
        li.classList.remove('visited', 'matched-path');
    });
}

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
function renderTree(node, depth = 0) {
    if (!node) return '';
    const idx = node.nodeIndex;
    const classNames = node.classes ? node.classes.join(' ') : '';
    const idAttr = node.id ? `#${node.id}` : '';
    const tagDisplay = `<span class="node-element" data-index="${idx}">
        <strong>${node.tag}</strong>${idAttr}${classNames ? ` .${classNames}` : ''}
    </span>`;
    let childrenHtml = '';
    if (node.children && node.children.length) {
        childrenHtml = '<ul>' + node.children.map(child => renderTree(child, depth + 1)).join('') + '</ul>';
    }
    return `<li>${tagDisplay}${childrenHtml}</li>`;
}

function renderLog(entries) {
    return entries.map(entry => {
        const badge = entry.matched ? '<span class="matched-badge">✓ MATCH</span>' : '';
        const classStr = entry.nodeClass?.join(' ') || '-';
        return `<li>
            <strong>Step ${entry.step}</strong> | 
            Tag: ${entry.nodeTag} | 
            ID: ${entry.nodeId || '-'} | 
            Class: ${classStr} | 
            Depth: ${entry.depth} 
            ${badge}
        </li>`;
    }).join('');
}

// ========== ANIMASI ==========
async function animateTraversal(logs) {
    const delay = 300;
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
}

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
    highlightMatchedPaths(logs);
}

function highlightMatchedPaths(logs) {
    const matchedIndices = new Set();
    logs.forEach(entry => { if (entry.matched) matchedIndices.add(entry.nodeIndex); });
    if (matchedIndices.size === 0) return;
    matchedIndices.forEach(idx => {
        const nodeSpan = document.querySelector(`.node-element[data-index="${idx}"]`);
        if (nodeSpan) {
            let li = nodeSpan.closest('li');
            while (li && li.tagName === 'LI') {
                if (li.classList.contains('visited')) li.classList.add('matched-path');
                li = li.parentElement?.closest('li');
            }
        }
    });
}

// ========== API CALLS ==========
// tombol Build DOM Tree
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

// tombol Search (mode css)
searchBtn.addEventListener('click', async () => {
    if (currentMode !== 'css') {
        alert('Silakan pilih mode CSS Selector terlebih dahulu.');
        return;
    }
    const algorithm = document.getElementById('algorithm').value;
    const selector = document.getElementById('selector').value.trim();
    let topN = parseInt(document.getElementById('topN').value, 10);
    const isUrlMode = document.getElementById('inputSourceSelect').value === 'url';
    const enableAnimation = document.getElementById('enableAnimation').checked;

    if (isNaN(topN)) topN = -1;
    if (topN < -1) {
        alert("Top N tidak boleh kurang dari -1.");
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
        document.getElementById('visitCount').innerText = data.visitCount ?? '-';
        document.getElementById('maxDepth').innerText = data.maxDepth ?? '-';
        document.getElementById('timeMs').innerText = data.elapsedTime ?? '-';
        document.getElementById('matchedCount').innerText = data.matchedCount ?? 0;
        if (data.tree) {
            document.getElementById('treeContainer').innerHTML = '<ul>' + renderTree(data.tree) + '</ul>';
        } else {
            document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Tidak ada pohon DOM</p>';
        }

        resetLCA();
        if (data.traversalLog && data.traversalLog.length) {
            document.getElementById('logContainer').innerHTML = renderLog(data.traversalLog);
            resetHighlights();
            if (enableAnimation) {
                await animateTraversal(data.traversalLog);
                highlightMatchedPaths(data.traversalLog);
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

// tombol Search LCA (mode lca)
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

// tombol Reset Pilihan (mode lca)
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