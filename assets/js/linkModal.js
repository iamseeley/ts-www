document.addEventListener('DOMContentLoaded', function() {
  const modalContainer = document.getElementById('modal-container');
  let currentModal = null;
  let hideTimeout = null;

  document.querySelectorAll('main .external-link').forEach(link => {
    link.addEventListener('mouseover', function() {
      clearTimeout(hideTimeout);
      showModal(this);
    });

    link.addEventListener('mouseout', function(event) {
      const toElement = event.relatedTarget || event.toElement;
      if (!toElement || !toElement.closest('.modal')) {
        hideTimeout = setTimeout(hideModal, 300);
      }
    });

    link.addEventListener('click', function(event) {
      event.preventDefault();
      clearTimeout(hideTimeout);
      showModal(this);
    });
  });

  function showModal(element) {
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
      <h2>${title}</h2>
      <p>${description}</p>
      <img src="${image}" alt="${title}">
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
    setTimeout(function() {
      modal.classList.add('show');
    }, 0);
  }

  function hideModal() {
    if (!currentModal) {
      return;
    }

    currentModal.classList.remove('show');
    setTimeout(function() {
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