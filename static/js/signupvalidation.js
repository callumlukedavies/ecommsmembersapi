

document.getElementById('form-signup').addEventListener('submit', function (event) {
    event.preventDefault();

    firstName = document.getElementById('signup-firstname');
    lastName = document.getElementById('signup-lastname');
    emailAddress = document.getElementById('signup-emailaddress');
    dateOfBirth = document.getElementById('signup-dateofbirth');
    password = document.getElementById('signup-password');
  
    firstNameErr = document.getElementById('firstname-error');
    lastNameErr = document.getElementById('lastname-error');
    emailAddressErr = document.getElementById('emailaddress-error');
    dateOfBirthErr = document.getElementById('dateofbirth-error');
    passwordErr = document.getElementById('password-error');
    generalErr = document.getElementById('general-error');

    firstNameErr.textContent = '';
    lastNameErr.textContent = '';
    emailAddressErr.textContent = '';
    dateOfBirthErr.textContent = '';
    passwordErr.textContent = '';

    const formData = new FormData(this);

    fetch('/membersapi/createuser', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
   
            return response.json().then(errorData => {
                if (errorData.errors) {
                   
                    if (errorData.errors.firstname) {
                        firstNameError.textContent = errorData.errors.firstname;
                    }
                    if (errorData.errors.lastname) {
                        lastNameErr.textContent = errorData.errors.lastname;
                    }
                    if (errorData.errors.email) {
                        emailAddressErr.textContent = errorData.errors.email;
                    }
                    if (errorData.errors.dateofbirth) {
                        dateOfBirthErr.textContent = errorData.errors.dateofbirth;
                    }
                    if (errorData.errors.password) {
                        passwordErr.textContent = errorData.errors.password;
                    }
                } else {
                    
                    generalErr.textContent = 'An unexpected error occurred. Please try again later.';
                }
            });
        } else {
            window.location.href = '/shopapi/';
            return response.json();
        }
    })
    .then(data => {
        if (data && data.success) {
            alert('Form submitted successfully!');
        }
    })
    .catch(error => {
        generalErr.textContent = 'An unexpected error occurred. Please try again later.';
    });
});
