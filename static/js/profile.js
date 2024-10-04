const addItemForm = document.getElementById('addItemForm');

// Clear the Add Item Form when submitted
addItemForm.addEventListener('submit', function(event){
    event.preventDefault();

    setTimeout(function() {
        addItemForm.reset();
    }, 100);
});