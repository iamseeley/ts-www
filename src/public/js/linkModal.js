document.addEventListener('DOMContentLoaded', function() {
  const modalContainer = document.getElementById('modal-container');
  let currentModal = null;
  let showTimeout = null;
  let hideTimeout = null;

  document.querySelectorAll('main .external-link').forEach(link => {
    link.addEventListener('mouseenter', function(event) {
      clearTimeout(hideTimeout);
      clearTimeout(showTimeout);
      showTimeout = setTimeout(() => showModal(this, event), 700);
    });

    link.addEventListener('mouseleave', function(event) {
      // Check if the mouse leaves the link but doesn't enter the modal
      if (!event.relatedTarget || !event.relatedTarget.closest('.modal')) {
        clearTimeout(showTimeout);
        hideTimeout = setTimeout(hideModal, 300);
      }
    });

    link.addEventListener('click', function(event) {
      event.preventDefault();
      clearTimeout(hideTimeout);
      clearTimeout(showTimeout);
      showModal(this, event);
    });
  });

  function showModal(element, event) {
    if (currentModal) {
      return;
    }
  
    const title = element.getAttribute('data-title');
    const description = element.getAttribute('data-description');
    const image = element.getAttribute('data-image');
    const url = element.href;
  
    const modal = document.createElement('div');
   
    modal.classList.add('modal');
    modal.innerHTML = `
      <p>${description}</p>
      <img src="${image}" alt="${title}" style="max-width: 100px; max-height: 100px;">
      <button class="proceed-btn">Proceed</button>
      <button class="close-btn">Close</button>
    `;
  
    modalContainer.innerHTML = '';
    modalContainer.appendChild(modal);
  
    modal.querySelector('.proceed-btn').onclick = function() {
      window.location.href = url;
    };
  
    modal.querySelector('.close-btn').addEventListener('click', hideModal);
  
    currentModal = modal;
    modalContainer.classList.add('show');
  
    // Get link position relative to viewport and add scroll offset to calculate absolute position
    const linkRect = element.getBoundingClientRect();
    const modalRect = modal.getBoundingClientRect();
    let topPosition = linkRect.bottom + window.scrollY + 6;
  
    let leftPosition = event.clientX - (modalRect.width / 2);
  
    // Adjust vertical position to ensure the modal is visible within the viewport
    if (topPosition + modalRect.height > window.scrollY + window.innerHeight) {
      topPosition = linkRect.top + window.scrollY - (modalRect.height + 8); // Place above the link if it doesn't fit below
    }
    if (leftPosition < 0) {
      leftPosition = 10; // Padding from the left edge of the viewport
    } else if (leftPosition + modalRect.width > window.innerWidth) {
      leftPosition = window.innerWidth - modalRect.width - 10; // Padding from the right edge of the viewport
    }
  
    modal.style.top = `${topPosition}px`;
    modal.style.left = `${leftPosition}px`;

    modal.addEventListener('mouseleave', function(modalEvent) {
      if (!modalEvent.relatedTarget || !modalEvent.relatedTarget.closest('.external-link')) {
        hideTimeout = setTimeout(hideModal, 300);
      }
    });
  
    setTimeout(() => modal.classList.add('show'), 0);
  }
  
  

  function hideModal() {
    if (!currentModal) {
      return;
    }

    currentModal.classList.remove('show');
    setTimeout(() => {
      modalContainer.classList.remove('show');
      modalContainer.innerHTML = '';
      currentModal = null;
    }, 300);
  }

  modalContainer.addEventListener('mouseover', function() {
    clearTimeout(hideTimeout);
  });

  modalContainer.addEventListener('mouseout', function() {
    hideTimeout = setTimeout(hideModal, 300);
  });
});
