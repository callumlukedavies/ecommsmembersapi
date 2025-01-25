document.getElementById('form-login').addEventListener('submit', function (event) {
    event.preventDefault();

    emailAddressErr = document.getElementById('emailaddress-error');
    passwordErr = document.getElementById('password-error');
    generalErr = document.getElementById('general-error');

    emailAddressErr.textContent = '';
    passwordErr.textContent = '';
    generalErr.textContent = '';

    const formData = new FormData(this);

    fetch('/membersapi/login', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(errorData => {

                if (errorData.errors) {
                    if (errorData.errors.email) {
                        emailAddressErr.textContent = errorData.errors.email;
                    }
                    if (errorData.errors.password) {
                        passwordErr.textContent = errorData.errors.password;
                    }
                } else {
                    generalErr.textContent = 'First Err: An unexpected error occurred. Please try again later.';
                }
            });
        } else {
            window.location.href = '/shopapi/profile';
            return response.json();
        }
    })
    .then(data => {
        if (data && data.success) {
            alert('Form submitted successfully!');
            console.log("Form submitted successfully");

            window.location.href = '/shopapi/profile'; // Redirect to the profile page
        }
    })
    .catch(error => {
        console.log(error);
    });
});