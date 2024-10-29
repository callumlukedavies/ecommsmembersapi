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

cartButtonClicked = false;
addToCartError = document.getElementById('add-to-cart-error');
addToCartButton = document.getElementById('add-to-cart-btn');
addToCartButton.addEventListener('click', function(event){
    if (cartButtonClicked) {
        console.log("Already clicked");
        return
    }

    cartButtonClicked = true;
    addToCartButton.innerHTML = 'Added to cart';
    const itemID = this.getAttribute('data-itemid');
    const sessionData = {
            itemID: itemID 
    }

    fetch('/shopapi/updatecart', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(sessionData)
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('UpdateCartHandler could not add item to cart.');
    })
    .then(data => {
        console.log('Success:', data); // Handle success (e.g., show message)
        if (data.error) {
            addToCartError.textContent = data.error;
        } else {
            addToCartError.textContent = data.message;
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});
