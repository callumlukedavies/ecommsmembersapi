document.getElementById('editprofile-firstname').addEventListener('submit', function (event) {
    event.preventDefault();

    firstName = document.getElementById('firstname-input');
    firstNameErr = document.getElementById('firstname-error');
    firstNameErr.textContent = '';

    const formData = new FormData(this);

    fetch('/membersapi/edit-user-firstname', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
   
            return response.json().then(errorData => {
                if (errorData.errorMessage) {
                    firstNameErr.textContent = errorData.errorMessage;
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

document.getElementById('editprofile-lastname').addEventListener('submit', function (event) {
    event.preventDefault();

    lastName = document.getElementById('lastname-input');
    lastNameErr = document.getElementById('lastname-error');
    lastNameErr.textContent = '';

    const formData = new FormData(this);

    fetch('/membersapi/edit-user-lastname', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
   
            return response.json().then(errorData => {
                if (errorData.errorMessage) {
                    lastNameErr.textContent = errorData.errorMessage;
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

document.getElementById('editprofile-email').addEventListener('submit', function (event) {
    event.preventDefault();

    emailAddress = document.getElementById('emailaddress-input');
    emailAddressErr = document.getElementById('emailaddress-error');


    emailAddressErr.textContent = '';

    const formData = new FormData(this);

    fetch('/membersapi/edit-user-emailaddress', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
   
            return response.json().then(errorData => {
                if (errorData.errorMessage) {
                    emailAddressErr.textContent = errorData.errorMessage;
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

document.getElementById('editprofile-dateofbirth').addEventListener('submit', function (event) {
    event.preventDefault();

    dateOfBirth = document.getElementById('dateofbirth-input');
    dateOfBirthErr = document.getElementById('dateofbirth-error');

    dateOfBirthErr.textContent = '';

    const formData = new FormData(this);

    fetch('/membersapi/edit-user-dateofbirth', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
   
            return response.json().then(errorData => {
                if (errorData.errorMessage) {
                    dateOfBirthErr.textContent = errorData.errorMessage;
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

document.getElementById('editprofile-password').addEventListener('submit', function (event) {
    event.preventDefault();

    password = document.getElementById('password-input');
    passwordErr = document.getElementById('password-error');
    passwordErr.textContent = '';

    const formData = new FormData(this);

    fetch('/membersapi/edit-user-password', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
   
            return response.json().then(errorData => {
                if (errorData.errorMessage) {
                    passwordErr.textContent = errorData.errorMessage;
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