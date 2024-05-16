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
  
    // Add touch event listeners to nav links
    navLinks.forEach(link => {
      link.addEventListener('touchstart', () => {
        link.classList.add('touch-active');
      });
      link.addEventListener('touchend', () => {
        setTimeout(() => link.classList.remove('touch-active'), 300); // Delay to see the active state
      });
    });
  });
  