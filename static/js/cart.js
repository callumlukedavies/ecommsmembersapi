const cartItems = document.querySelectorAll('.cart-item');
const removeButtons = document.querySelectorAll('.remove-btn');
const subtotalText = document.getElementById('cart-subtotal');
const deliveryText = document.getElementById('cart-delivery');
const totalText = document.getElementById('cart-total');
const itemNumberText = document.getElementById('cart-item-number');

var subtotal = 0;
var delivery = 5;
var total = 0;
var noOfItems = 0;

cartItems.forEach(cartItem => {
    const price = cartItem.getAttribute('data-itemprice'); 
    const parsedPrice = parseInt(price, 10);
    noOfItems++;

    subtotal += parsedPrice;

    cartItem.querySelector(".remove-btn").addEventListener('click', function(){
        
        const sessionData = {
            itemID:  cartItem.getAttribute('data-itemid')
        }

        fetch('/shopapi/removecartitem', {
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
            throw new Error('removecartitem could not remove item from cart.');
        })
        .then(data => {
            console.log('Success:', data);

            subtotal -= parsedPrice;
            total = subtotal + delivery;
            noOfItems--;

            subtotalText.innerHTML = "Subtotal: £" + subtotal;
            totalText.innerHTML = "Total: £" + total;
            itemNumberText.innerHTML = "Number of items: " + noOfItems;
        })
        .catch(error => {
            console.error('Error:', error);
        });

        console.log("removing cartitem");
        cartItem.remove();
    });
});

total = subtotal + delivery;

subtotalText.innerHTML = "Subtotal: £" + subtotal;
deliveryText.innerHTML = "Delivery: £5";
totalText.innerHTML = "Total: £" + total;
itemNumberText.innerHTML = "Number of items: " + noOfItems;

checkoutCartButton = document.getElementById('checkout-btn');
checkoutCartButton.addEventListener('click', function(){
    if (noOfItems == 0) {
        return
    }

    items = [];

    cartItems.forEach(cartItem => {
        itemID = cartItem.getAttribute('data-itemid');
        console.log(itemID);
        items.push(itemID);
    });

    const sessionData = {
            itemIDs: items
    }

    fetch('/shopapi/checkout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(sessionData)
    })
    .then(response => {
        if (response.ok) {
            window.location.href = "/shopapi/";
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