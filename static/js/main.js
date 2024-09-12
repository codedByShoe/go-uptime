// import bootstrap js
import "bootstrap/js/dist/button";
import "bootstrap/js/dist/collapse";
import "bootstrap/js/dist/dropdown";
import Modal from "bootstrap/js/dist/modal";
import Alert from "bootstrap/js/dist/alert";

// hide back-btn if route is /
const route = window.location.pathname;
if (route === "/") {
  const backBtn = document.getElementById("back-btn");
  backBtn.classList.add("d-none");
}

// get add-site-modal
const addSiteModal = new Modal(document.getElementById("add-site-modal"));
// add-site-btn
const addSiteBtn = document.getElementById("add-site-btn");
addSiteBtn.addEventListener("click", function () {
  addSiteModal.show();
});

// get error-alert
const errorAlert = new Alert(document.getElementById("error-alert"));
// get error-alert-close
const errorAlertClose = document.getElementById("error-alert-close");
errorAlertClose.addEventListener("click", function () {
  errorAlert.close();
});
