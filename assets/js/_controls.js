/**
 * All general initializations are done in this file.
 * Controls such as menu, multi dropdown needs be initialized
 * as how semantic-ui suggests. This is irrespective of element id
 * So initialized here and this file gets loaded in every page and
 * does all the necessary initializations.
 */

// Global object can be used in all files.
window.QaNet = window.QaNet || {};

// Menu dropdown initialization.
const initDropdown = function () {
  $('#menu').dropdown();

  $('.multidd').dropdown({
    allowAdditions: true
  });
};

const closeMessage = function () {
  $('.message .close').on('click', function () {
    $(this)
      .closest('.message')
      .transition('fade');
  });
};

const closeShorErrorMessage = function () {
  $('.short-error-message').on('click', function () {
    $(this).addClass('d-n');
  });
};

export default function init() {
  initDropdown();
  closeMessage();
  closeShorErrorMessage();
}
