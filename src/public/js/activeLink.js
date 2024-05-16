document.addEventListener('DOMContentLoaded', () => {
    
    const navLinks = document.querySelectorAll('nav a');
    

    function updateActiveLink() {
  
      const path = window.location.pathname;
  

      navLinks.forEach(link => {

        if (link.getAttribute('href') === path) {
          link.classList.add('active');
        } else {
          link.classList.remove('active');
        }
      });
    }
  
  
    updateActiveLink();
    
    window.addEventListener('popstate', updateActiveLink);
  });
  