
import './style.css';

async function selectDataFile() {
  const filePath = await window.backend.App.SelectFile("Select Data.csv File");
  document.getElementById("dataFile").value = filePath;
}

async function selectStudentCodesFile() {
  const filePath = await window.backend.App.SelectFile("Select Student Codes.txt File");
  document.getElementById("studentCodesFile").value = filePath;
}

async function selectStaffCodesFile() {
  const filePath = await window.backend.App.SelectFile("Select Staff Codes.txt File");
  document.getElementById("staffCodesFile").value = filePath;
}

async function selectIDTemplateFile() {
  const filePath = await window.backend.App.SelectFile("Select ID Template File");
  document.getElementById("idTemplateFile").value = filePath;
}

async function generateIDCards() {
  const dataFile = document.getElementById("dataFile").value;
  const studentCodesFile = document.getElementById("studentCodesFile").value;
  const staffCodesFile = document.getElementById("staffCodesFile").value;
  const idTemplateFile = document.getElementById("idTemplateFile").value;

  if (!dataFile || !studentCodesFile || !staffCodesFile || !idTemplateFile) {
    alert("Please select all required files before generating ID cards.");
    return;
  }

  const result = await window.backend.App.GenerateIDCards({
    dataFile,
    studentCodesFile,
    staffCodesFile,
    idTemplateFile,
  });

  alert(result);
}
