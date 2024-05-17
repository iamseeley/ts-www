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
      showTimeout = setTimeout(() => showModal(this, event), 600);
    });

    link.addEventListener('mouseleave', function(event) {
      if (!event.relatedTarget || !event.relatedTarget.closest('.modal')) {
        clearTimeout(showTimeout);
        hideTimeout = setTimeout(hideModal, 400);
      }
    });

    link.addEventListener('click', function(event) {
      if (!currentModal) {
        event.preventDefault();
        clearTimeout(hideTimeout);
        clearTimeout(showTimeout);
        showModal(this, event);
      }
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
    const isExternal = element.target === '_blank';

    const modal = document.createElement('a');
    modal.classList.add('modal');
    modal.href = url;
    modal.target = isExternal ? '_blank' : '_self';
    modal.innerHTML = `
      <div>
        <p>${description}</p>
        <img src="${image}" alt="${title}" style="max-width: 100px; max-height: 100px;">
      </div>
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

    const containerRect = modalContainer.getBoundingClientRect();
    if (leftPosition < containerRect.left) {
      leftPosition = containerRect.left + 10;
    } else if (leftPosition + modalRect.width > containerRect.right) {
      leftPosition = containerRect.right - modalRect.width - 10;
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
