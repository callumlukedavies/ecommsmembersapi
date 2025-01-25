const categoryInput = document.getElementById('category-input');
const sizeInput = document.getElementById('size-input');
const deleteItemBtn = document.getElementById('delete-item-btn');
const editItemForm = document.getElementById('edit-item-form');

// sizeInput.innerHTML = '<option value="">Please select a category before choose size</option>';

// editItemForm.addEventListener('submit', function(){
//     const formData = new FormData(this);
    
//     fetch('/shopapi/', {
//         method: 'POST',
//         body: formData
//     })
//     .then(response => {
//         if (!response.ok) {
//             return response.json().then(errorData => {
//                 console.log("response not ok");
//                 if (errorData.errors) {
//                     console.log("errordata.errors has triggered");
//                 } else {
//                     generalErr.textContent = 'First Err: An unexpected error occurred. Please try again later.';
//                 }
//             });
//         } else {
//             console.log("response ok");
//             window.location.href = '/shopapi/';
//             return response.json();
//         }
//     })
//     .then(data => {
//         if (data && data.success) {
//             alert('Form submitted successfully!');
//             console.log("Form submitted successfully");
//         }
//         console.log("Redirecting to home");
//         window.location.href = '/shopapi/'; // Redirect to the profile page

//     })
//     .catch(error => {
//         console.log(error);
//     });
// });

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