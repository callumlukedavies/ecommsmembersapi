const itemCards = document.querySelectorAll('.item-card');
const deleteButtons = document.querySelectorAll(".delete-item-btn");

itemCards.forEach(itemCard => {
    itemCard.querySelector(".delete-item-btn").addEventListener('click', function(){
        const itemId = this.getAttribute("data-itemid");

        fetch(`/shopapi/${itemId}`, {
            method: "DELETE",
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Failed to delete the item");
            }
            return response.json();
        })
        .then(data => {
            alert(data.message);
            itemCard.remove();
        })
        .catch(error => {
            console.error("Error:", error);
            alert("Failed to delete the item");
        });
    });
});