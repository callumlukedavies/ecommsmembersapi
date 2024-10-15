const addItemForm = document.getElementById('addItemForm');
const categoryInput = document.getElementById('category-input');
const sizeInput = document.getElementById('size-input');

sizeInput.innerHTML = '<option value="">Please select a category before choose size</option>';

// Clear the Add Item Form when submitted
addItemForm.addEventListener('submit', function(event){
    // event.preventDefault();

    // setTimeout(function() {
    //     addItemForm.reset();
    // }, 100);

    console.log("Form being submitted")
});

categoryInput.addEventListener('change', function(event){
    console.log("Category changed!");

    const category = this.value;
      
    sizeInput.innerHTML = '';

    let sizes = [];

    if (category === 'Bottoms') {
        sizes = [
            { value: '28', text: '28' },
            { value: '30', text: '30' },
            { value: '32', text: '32' },
            { value: '34', text: '34' },
            { value: '36', text: '36' }
        ];
    } else if (category === 'Shoe') {
        sizes = [
            { value: '3', text: '3' },
            { value: '4', text: '4' },
            { value: '5', text: '5' },
            { value: '6', text: '6' },
            { value: '7', text: '7' },
            { value: '8', text: '8' },
            { value: '9', text: '9' },
            { value: '10', text: '10' },
            { value: '11', text: '11' },
            { value: '12', text: '12' }
        ];
    } else {
        sizes = [
            { value: 'XS', text: 'XS' },
            { value: 'S', text: 'S' },
            { value: 'M', text: 'M' },
            { value: 'L', text: 'L' },
            { value: 'XL', text: 'XL' },
        ];
    }

    sizes.forEach(function(option) {
        const newOption = document.createElement('option');
        newOption.value = option.value;
        newOption.textContent = option.text;
        sizeInput.appendChild(newOption);
      });
});