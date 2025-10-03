// Gin Static Files Example JavaScript

document.addEventListener('DOMContentLoaded', function() {
    console.log('Gin Static Files Example loaded');
    loadFiles();

    // 파일 업로드 폼
    const uploadForm = document.getElementById('uploadForm');
    uploadForm.addEventListener('submit', handleUpload);
});

// 파일 업로드 처리
async function handleUpload(e) {
    e.preventDefault();

    const fileInput = document.getElementById('fileInput');
    const file = fileInput.files[0];
    const uploadResult = document.getElementById('uploadResult');

    if (!file) {
        uploadResult.className = 'error';
        uploadResult.textContent = '파일을 선택해주세요';
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
        const response = await fetch('/upload', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (response.ok) {
            uploadResult.className = 'success';
            uploadResult.textContent = `✅ 업로드 성공: ${data.filename} (${formatFileSize(data.size)})`;
            fileInput.value = '';
            loadFiles();
        } else {
            uploadResult.className = 'error';
            uploadResult.textContent = `❌ 업로드 실패: ${data.error}`;
        }
    } catch (error) {
        uploadResult.className = 'error';
        uploadResult.textContent = `❌ 오류 발생: ${error.message}`;
    }
}

// 파일 목록 로드
async function loadFiles() {
    const filesList = document.getElementById('filesList');
    filesList.innerHTML = '<div class="loading">파일 목록 로딩 중...</div>';

    try {
        const response = await fetch('/api/files');
        const data = await response.json();

        if (data.files && data.files.length > 0) {
            filesList.innerHTML = data.files.map(file => `
                <div class="file-item">
                    <div class="file-info">
                        <div class="file-name">📄 ${file.name}</div>
                        <div class="file-size">${formatFileSize(file.size)}</div>
                    </div>
                    <div class="file-actions">
                        <a href="${file.url}" download class="download-btn">
                            <button>⬇️ 다운로드</button>
                        </a>
                        <button class="delete-btn" onclick="deleteFile('${file.name}')">
                            🗑️ 삭제
                        </button>
                    </div>
                </div>
            `).join('');
        } else {
            filesList.innerHTML = '<div class="empty-state">업로드된 파일이 없습니다</div>';
        }
    } catch (error) {
        filesList.innerHTML = `<div class="error">파일 목록 로드 실패: ${error.message}</div>`;
    }
}

// 파일 삭제
async function deleteFile(filename) {
    if (!confirm(`'${filename}' 파일을 삭제하시겠습니까?`)) {
        return;
    }

    try {
        const response = await fetch(`/api/files/${filename}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            alert('파일이 삭제되었습니다');
            loadFiles();
        } else {
            const data = await response.json();
            alert(`삭제 실패: ${data.error}`);
        }
    } catch (error) {
        alert(`오류 발생: ${error.message}`);
    }
}

// 파일 크기 포맷팅
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 이미지 미리보기 (추가 기능)
function previewImage(file) {
    const reader = new FileReader();
    reader.onload = function(e) {
        const preview = document.createElement('img');
        preview.src = e.target.result;
        preview.style.maxWidth = '200px';
        preview.style.marginTop = '10px';
        document.getElementById('uploadResult').appendChild(preview);
    };

    if (file && file.type.startsWith('image/')) {
        reader.readAsDataURL(file);
    }
}

// 드래그 앤 드롭 지원
const fileInput = document.getElementById('fileInput');
if (fileInput) {
    const uploadArea = fileInput.parentElement;

    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        uploadArea.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        uploadArea.addEventListener(eventName, highlight, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        uploadArea.addEventListener(eventName, unhighlight, false);
    });

    function highlight(e) {
        uploadArea.classList.add('highlight');
    }

    function unhighlight(e) {
        uploadArea.classList.remove('highlight');
    }

    uploadArea.addEventListener('drop', handleDrop, false);

    function handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;

        if (files.length > 0) {
            fileInput.files = files;
            handleUpload(e);
        }
    }
}