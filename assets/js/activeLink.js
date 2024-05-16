document.addEventListener('DOMContentLoaded', () => {
    // Get all nav links
    const navLinks = document.querySelectorAll('nav a');
  
    // Function to update the active link
    function updateActiveLink() {
      // Get the current path
      const path = window.location.pathname;
  
      // Iterate over all nav links
      navLinks.forEach(link => {
        // Check if the link href matches the current path
        if (link.getAttribute('href') === path) {
          link.classList.add('active');
        } else {
          link.classList.remove('active');
        }
      });
    }
  
    // Call the function to set the active link initially
    updateActiveLink();
  
    // Optional: Add a listener for navigation changes (if using a router)
    window.addEventListener('popstate', updateActiveLink);
  });