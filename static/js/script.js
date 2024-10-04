
// Get elements
const loginBtn = document.getElementById('loginBtn');
const signupBtn = document.getElementById('signupBtn');
const closeBtn = document.getElementById('closeBtn');
const popup = document.getElementById('popup');
const overlay = document.getElementById('overlay');
const submitCredsBtn = document.getElementById('submitCredentialsBtn');

// Show popup when login button is clicked
loginBtn.addEventListener('click', function(e) {
    e.preventDefault();
    popup.style.display = 'block';
    overlay.style.display = 'block';
});

// Hide popup when close button is clicked
closeBtn.addEventListener('click', function() {
    popup.style.display = 'none';
    overlay.style.display = 'none';
});

// Hide popup when clicking outside of the popup
overlay.addEventListener('click', function() {
    popup.style.display = 'none';
    overlay.style.display = 'none';
});


function toggleEditForm(id) {
    const form = document.getElementById(`editPrice${id}`);
    const isHidden = form.style.display === 'none' || form.style.display === '';
    form.style.display = isHidden ? 'block' : 'none';
};
