// Tab switching
function switchTab(type) {
    const urlContainer = document.getElementById('input-url-container');
    const rawContainer = document.getElementById('input-raw-container');
    const btnUrl = document.getElementById('btn-url');
    const btnRaw = document.getElementById('btn-raw');

    if (type === 'url') {
        urlContainer.classList.remove('hidden');
        rawContainer.classList.add('hidden');
        btnUrl.classList.add('active');
        btnRaw.classList.remove('active');
    } else {
        urlContainer.classList.add('hidden');
        rawContainer.classList.remove('hidden');
        btnUrl.classList.remove('active');
        btnRaw.classList.add('active');
    }
}

// ========== INTEGRASI BACKEND ==========
document.getElementById('searchBtn').addEventListener('click', async () => {
    const algorithm = document.getElementById('algorithm').value;
    const selector = document.getElementById('selector').value.trim();
    let topN = parseInt(document.getElementById('topN').value, 10);
    const isUrlMode = document.getElementById('btn-url').classList.contains('active');

    // Validasi TopN
    if (isNaN(topN)) topN = -1;
    if (topN < -1) {
        alert("Top N tidak boleh kurang dari -1. -1 berarti semua hasil.");
        document.getElementById('topN').value = -1;
        topN = -1;
    }

    // Validasi selector tidak kosong
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

    // Tampilkan loading
    document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Memproses...</p>';
    document.getElementById('logContainer').innerHTML = '<li class="hint">Memuat log...</li>';
    
    try {
        const response = await fetch('http://localhost:8080/api/search', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        const data = await response.json();

        if (!data.success) {
            alert('Error: ' + (data.message || 'Unknown error'));
            return;
        }

        // Update statistik
        document.getElementById('visitCount').innerText = data.visitCount ?? '-';
        document.getElementById('maxDepth').innerText = data.maxDepth ?? '-';
        document.getElementById('timeMs').innerText = data.elapsedTime ?? '-';
        document.getElementById('matchedCount').innerText = data.matchedCount ?? 0;

        // Render pohon DOM
        if (data.tree) {
            document.getElementById('treeContainer').innerHTML = '<ul>' + renderTree(data.tree) + '</ul>';
        } else {
            document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Tidak ada pohon DOM</p>';
        }

        // Tampilkan log traversal
        if (data.traversalLog && data.traversalLog.length > 0) {
            document.getElementById('logContainer').innerHTML = renderLog(data.traversalLog);
            resetHighlights();
            // Jalankan animasi traversal
            await animateTraversal(data.traversalLog);
            // Setelah animasi selesai, highlight jalur dari root ke setiap node yang matched
            highlightMatchedPaths(data.traversalLog);
        } else {
            document.getElementById('logContainer').innerHTML = '<li class="hint">Tidak ada log traversal.</li>';
        }
    } catch (err) {
        console.error(err);
        alert('Gagal menghubungi backend: ' + err.message);
        document.getElementById('treeContainer').innerHTML = '<p class="placeholder-text">Error koneksi ke server</p>';
    }
});

// Render tree dengan data-index
function renderTree(node, depth = 0) {
    if (!node) return '';
    const idx = node.nodeIndex;
    const classNames = node.classes ? node.classes.join(' ') : '';
    const idAttr = node.id ? `#${node.id}` : '';
    const tagDisplay = `<span class="node-element" data-index="${idx}">
        <strong>${node.tag}</strong>${idAttr}${classNames ? ` .${classNames}` : ''}
    </span>`;
    let childrenHtml = '';
    if (node.children && node.children.length > 0) {
        childrenHtml = '<ul>' + node.children.map(child => renderTree(child, depth + 1)).join('') + '</ul>';
    }
    return `<li>${tagDisplay}${childrenHtml}</li>`;
}

// Render log traversal
function renderLog(entries) {
    return entries.map(entry => {
        const matchedBadge = entry.matched ? '<span class="matched-badge">✓ MATCH</span>' : '';
        const classStr = entry.nodeClass?.join(' ') || '-';
        return `<li>
            <strong>Step ${entry.step}</strong> | 
            Tag: ${entry.nodeTag} | 
            ID: ${entry.nodeId || '-'} | 
            Class: ${classStr} | 
            Depth: ${entry.depth} 
            ${matchedBadge}
        </li>`;
    }).join('');
}

// Animasi traversal langkah per langkah
async function animateTraversal(logs) {
    const delay = 300;
    for (let i = 0; i < logs.length; i++) {
        const step = logs[i];
        
        // Highlight di panel log
        const logEntries = document.querySelectorAll('#logContainer li');
        if (logEntries[i]) {
            logEntries[i].classList.add('current-step');
        }

        // Highlight node di pohon
        const targetNode = document.querySelector(`.node-element[data-index="${step.nodeIndex}"]`);
        if (targetNode) {
            const parentLi = targetNode.closest('li');
            targetNode.classList.add('visited');
            if (parentLi) parentLi.classList.add('visited');
            if (step.matched) {
                targetNode.classList.add('matched');
            }
            targetNode.classList.add('current-step');
        }

        await new Promise(resolve => setTimeout(resolve, delay));
        
        if (targetNode) targetNode.classList.remove('current-step');
        if (logEntries[i]) logEntries[i].classList.remove('current-step');
    }
}

// Highlight jalur dari root ke setiap node yang ditemukan (matched)
// Hanya untuk ancestor yang sudah memiliki class 'visited' (sudah dikunjungi)
function highlightMatchedPaths(traversalLog) {
    if (!traversalLog || traversalLog.length === 0) return;
    
    // Kumpulkan nodeIndex yang matched
    const matchedIndices = new Set();
    traversalLog.forEach(entry => {
        if (entry.matched) {
            matchedIndices.add(entry.nodeIndex);
        }
    });
    if (matchedIndices.size === 0) return;

    matchedIndices.forEach(idx => {
        const nodeSpan = document.querySelector(`.node-element[data-index="${idx}"]`);
        if (nodeSpan) {
            const li = nodeSpan.closest('li');
            if (li) {
                // Naik ke ancestor, tapi hanya yang sudah memiliki class 'visited'
                let current = li;
                while (current && current.tagName === 'LI') {
                    if (current.classList.contains('visited')) {
                        current.classList.add('matched-path');
                    }
                    current = current.parentElement?.closest('li');
                }
            }
        }
    });
}

// Reset semua highlight sebelum animasi baru
function resetHighlights() {
    document.querySelectorAll('.node-element').forEach(el => {
        el.classList.remove('visited', 'matched', 'current-step');
    });
    document.querySelectorAll('.canvas-container li').forEach(li => {
        li.classList.remove('visited', 'matched-path');
    });
}