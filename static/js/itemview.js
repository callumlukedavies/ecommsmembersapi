$(document).ready(function() {
    var carouselWidth = $('.bottom-carousel-inner')[0].scrollWidth;
    var itemWidth = $('.bottom-carousel-item').outerWidth();

    var scrollPosition = 0;

    $('.carousel-control-next').on('click', function() {
        if (scrollPosition < (carouselWidth - (itemWidth * 5))) {
            scrollPosition = scrollPosition + itemWidth + 10;
            $('.bottom-carousel-inner').animate({ scrollLeft: scrollPosition }, 600);
        }
    });

    $('.carousel-control-prev').on('click', function() {
        if (scrollPosition > 0) {
            scrollPosition = scrollPosition - itemWidth -10;
            $('.bottom-carousel-inner').animate({ scrollLeft: scrollPosition }, 600);
        }
    });
});
