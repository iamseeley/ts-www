document.addEventListener('DOMContentLoaded', function() {
  const modalContainer = document.getElementById('modal-container');
  let currentModal = null;
  let showTimeout = null;
  let hideTimeout = null;

  document.querySelectorAll('main .external-link').forEach(link => {
    link.addEventListener('mouseenter', function(event) {
      clearTimeout(hideTimeout);
      clearTimeout(showTimeout);
      showTimeout = setTimeout(() => showModal(this, event), 600);
    });

    link.addEventListener('mouseleave', function(event) {
      if (!event.relatedTarget || !event.relatedTarget.closest('.modal')) {
        clearTimeout(showTimeout);
        hideTimeout = setTimeout(hideModal, 400);
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
    `;

    modalContainer.innerHTML = '';
    modalContainer.appendChild(modal);


    currentModal = modal;
    modalContainer.classList.add('show');

    const linkRect = element.getBoundingClientRect();
    const modalRect = modal.getBoundingClientRect();
    let topPosition = linkRect.bottom + window.scrollY + 6;
    let leftPosition = event.clientX - (modalRect.width / 2);

    if (topPosition + modalRect.height > window.scrollY + window.innerHeight) {
      topPosition = linkRect.top + window.scrollY - (modalRect.height + 8);
    }
    if (leftPosition < 0) {
      leftPosition = 10;
    } else if (leftPosition + modalRect.width > window.innerWidth) {
      leftPosition = window.innerWidth - modalRect.width - 10;
    }

    modal.style.top = `${topPosition}px`;
    modal.style.left = `${leftPosition}px`;

    setTimeout(() => modal.classList.add('show'), 0);

    document.addEventListener('touchstart', handleOutsideTouch, true);
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
      document.removeEventListener('touchstart', handleOutsideTouch, true);
    }, 300);
  }

  function handleOutsideTouch(event) {
    if (currentModal && !modalContainer.contains(event.target) && !event.target.closest('.external-link') && !event.target.closest('.modal')) {
      hideModal();
    }
  }
});
