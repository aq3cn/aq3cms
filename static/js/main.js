// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    console.log('页面加载完成');
    
    // 初始化轮播图
    initSlider();
    
    // 初始化搜索框
    initSearch();
});

// 初始化轮播图
function initSlider() {
    const slides = document.querySelectorAll('.slide');
    if (slides.length <= 1) return;
    
    let currentSlide = 0;
    
    // 隐藏所有幻灯片，只显示当前幻灯片
    function showSlide(index) {
        slides.forEach((slide, i) => {
            if (i === index) {
                slide.style.display = 'block';
            } else {
                slide.style.display = 'none';
            }
        });
    }
    
    // 显示下一张幻灯片
    function nextSlide() {
        currentSlide = (currentSlide + 1) % slides.length;
        showSlide(currentSlide);
    }
    
    // 初始化显示第一张幻灯片
    showSlide(0);
    
    // 设置定时器，自动切换幻灯片
    setInterval(nextSlide, 5000);
}

// 初始化搜索框
function initSearch() {
    const searchForm = document.querySelector('.search form');
    if (!searchForm) return;
    
    searchForm.addEventListener('submit', function(e) {
        const input = this.querySelector('input');
        if (!input.value.trim()) {
            e.preventDefault();
            alert('请输入搜索关键词');
        }
    });
}
