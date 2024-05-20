document.addEventListener('DOMContentLoaded', function() {
    const modalContainer = document.getElementById('modal-container');
    let currentModal = null;
    let showTimeout = null;
    let hideTimeout = null;
  
    document.querySelectorAll('main a').forEach(link => {
      link.addEventListener('mouseenter', function(event) {
        clearTimeout(hideTimeout);
        clearTimeout(showTimeout);
        if (currentModal) {
          hideModal();
        }
        const description = this.getAttribute('data-description');
        const image = this.getAttribute('data-image');
        if (description || image) {
          showTimeout = setTimeout(() => showModal(this, event), 600);
        }
      });
  
      link.addEventListener('mouseleave', function(event) {
        if (!event.relatedTarget || !event.relatedTarget.closest('.modal')) {
          clearTimeout(showTimeout);
          hideTimeout = setTimeout(hideModal, 400);
        }
      });
  
      link.addEventListener('click', function(event) {
        if (!currentModal) {
          const description = this.getAttribute('data-description');
          const image = this.getAttribute('data-image');
          if (description || image) {
            event.preventDefault();
            clearTimeout(hideTimeout);
            clearTimeout(showTimeout);
            showModal(this, event);
          }
        }
      });
    });
  
    function showModal(element, event) {
      const description = element.getAttribute('data-description');
      const image = element.getAttribute('data-image');
  
      if (!description && !image) {
        return;
      }
  
      if (currentModal) {
        return;
      }
  
      const title = element.getAttribute('data-title');
      const url = element.href;
      const isExternal = element.target === '_blank';
  
      const modal = document.createElement('a');
      modal.classList.add('modal');
      modal.href = url;
      modal.target = isExternal ? '_blank' : '_self';
  
      let modalContent = '';
      if (description) {
        modalContent += `<p>${description}</p>`;
      }
      if (image) {
        modalContent += `<img src="${image}" alt="${title}" style="max-width: 100px; max-height: 100px;">`;
      }
  
      modal.innerHTML = `<div>${modalContent}</div>`;
  
      modalContainer.innerHTML = '';
      modalContainer.appendChild(modal);
      currentModal = modal;
      modalContainer.classList.add('show');
  
      const linkRect = element.getBoundingClientRect();
      const modalRect = modal.getBoundingClientRect();
  
      let topPosition = linkRect.bottom + window.scrollY + 6;
      let leftPosition = linkRect.left + window.scrollX - (modalRect.width / 2) + (linkRect.width / 2);
  
      // Ensure the modal doesn't overflow the viewport boundaries
      if (topPosition + modalRect.height > window.scrollY + window.innerHeight) {
        topPosition = linkRect.top + window.scrollY - modalRect.height - 6;
      }
  
      if (leftPosition < 0) {
        leftPosition = 10;
      } else if (leftPosition + modalRect.width > window.scrollX + window.innerWidth) {
        leftPosition = window.scrollX + window.innerWidth - modalRect.width - 10;
      }

     // Adjust left position for larger screens
     const screenWidth = window.innerWidth;
     if (screenWidth > 850) {  // Assuming 1024px as the threshold for larger screens
       leftPosition += 100; // Move slightly to the right of the link
     }
    
      modal.style.top = `${topPosition}px`;
      modal.style.left = `${leftPosition}px`;
  
      setTimeout(() => modal.classList.add('show'), 0);
  
      document.addEventListener('touchstart', handleOutsideTouch, true);
      modal.addEventListener('mouseleave', handleModalMouseLeave);
      modal.addEventListener('mouseenter', handleModalMouseEnter);
    }
  
    function hideModal() {
      if (!currentModal) {
        return;
      }
  
      currentModal.classList.remove('show');
      setTimeout(() => {
        if (currentModal) {
          modalContainer.classList.remove('show');
          modalContainer.innerHTML = '';
          currentModal.removeEventListener('mouseleave', handleModalMouseLeave);
          currentModal.removeEventListener('mouseenter', handleModalMouseEnter);
          currentModal = null;
        }
        document.removeEventListener('touchstart', handleOutsideTouch, true);
      }, 300);
    }
  
    function handleOutsideTouch(event) {
      if (currentModal && !modalContainer.contains(event.target) && !event.target.closest('.external-link') && !event.target.closest('.modal')) {
        hideModal();
      }
    }
  
    function handleModalMouseLeave() {
      hideTimeout = setTimeout(hideModal, 400);
    }
  
    function handleModalMouseEnter() {
      clearTimeout(hideTimeout);
    }
  });
  