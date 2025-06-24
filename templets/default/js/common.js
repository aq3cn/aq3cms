/**
 * aq3cms前台公共JavaScript
 */
$(function() {
    // 设置当前导航高亮
    setCurrentNav();
    
    // 返回顶部按钮
    initBackToTop();
    
    // 图片懒加载
    initLazyLoad();
    
    // 响应式表格
    initResponsiveTables();
    
    // 文章内容图片点击放大
    initImageZoom();
});

/**
 * 设置当前导航高亮
 */
function setCurrentNav() {
    var pathname = window.location.pathname;
    
    // 首页
    if (pathname === '/' || pathname === '/index.html') {
        $('.navbar-nav li:first').addClass('active');
        return;
    }
    
    // 其他页面
    $('.navbar-nav li').each(function() {
        var link = $(this).find('a').attr('href');
        if (pathname.indexOf(link) === 0 && link !== '/') {
            $(this).addClass('active').siblings().removeClass('active');
        }
    });
}

/**
 * 初始化返回顶部按钮
 */
function initBackToTop() {
    // 添加返回顶部按钮
    $('body').append('<a href="javascript:;" class="back-to-top" title="返回顶部"><i class="glyphicon glyphicon-chevron-up"></i></a>');
    
    var $backToTop = $('.back-to-top');
    
    // 滚动时显示/隐藏按钮
    $(window).scroll(function() {
        if ($(this).scrollTop() > 300) {
            $backToTop.fadeIn();
        } else {
            $backToTop.fadeOut();
        }
    });
    
    // 点击返回顶部
    $backToTop.click(function() {
        $('html, body').animate({
            scrollTop: 0
        }, 800);
        return false;
    });
    
    // 添加样式
    var style = `
        .back-to-top {
            position: fixed;
            right: 20px;
            bottom: 20px;
            width: 40px;
            height: 40px;
            line-height: 40px;
            text-align: center;
            background-color: rgba(0, 0, 0, 0.5);
            color: #fff;
            border-radius: 4px;
            display: none;
            z-index: 9999;
        }
        .back-to-top:hover {
            background-color: rgba(0, 0, 0, 0.7);
            color: #fff;
            text-decoration: none;
        }
    `;
    
    $('head').append('<style>' + style + '</style>');
}

/**
 * 初始化图片懒加载
 */
function initLazyLoad() {
    // 如果没有引入懒加载插件，则不执行
    if (typeof $.fn.lazyload !== 'function') {
        return;
    }
    
    // 图片懒加载
    $('img.lazy').lazyload({
        effect: 'fadeIn',
        threshold: 200
    });
}

/**
 * 初始化响应式表格
 */
function initResponsiveTables() {
    // 给文章内容中的表格添加响应式样式
    $('.article-content table').addClass('table table-bordered table-striped table-hover');
    
    // 添加响应式容器
    $('.article-content table').wrap('<div class="table-responsive"></div>');
}

/**
 * 初始化图片点击放大
 */
function initImageZoom() {
    // 如果没有引入图片查看器插件，则不执行
    if (typeof $.fn.viewer !== 'function') {
        return;
    }
    
    // 初始化图片查看器
    $('.article-content').viewer({
        url: 'data-original',
        navbar: false,
        title: false,
        toolbar: {
            zoomIn: true,
            zoomOut: true,
            oneToOne: true,
            reset: true,
            prev: true,
            next: true,
            rotateLeft: true,
            rotateRight: true
        }
    });
}
