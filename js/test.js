function sms() {
  var xmlhttp = new XMLHttpRequest(); // IE7+, Firefox, Chrome, Opera, Safari.
  xmlhttp.open('POST', '/sms', true);
  xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xmlhttp.send();   // add paarms if needed
}
