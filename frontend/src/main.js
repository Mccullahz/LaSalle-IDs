import './style.css';

import * as App from '../wailsjs/go/main/App';
import * as WailsRuntime from '../wailsjs/runtime';

WailsRuntime.EventsOn("dataFileSelected", (filePath) => {
  document.getElementById("dataFile").value = filePath;
});

WailsRuntime.EventsOn("idTemplateFileSelected", (filePath) => {
  document.getElementById("idTemplateFile").value = filePath;
});

WailsRuntime.EventsOn("outputDirectorySelected", (directoryPath) => {
  document.getElementById("outputDirectory").value = directoryPath;
});

function selectDataFile() {
  App.SelectDataFile();
}

function selectIDTemplateFile() {
  App.SelectIDTemplateFile();
}

function selectOutputDirectory() {
  App.SelectOutputDirectory();
}

async function generateIDCards() {
  const result = await App.GenerateIDCards();
  alert(result);
}

document.getElementById('selectDataFileButton').addEventListener('click', selectDataFile);
document.getElementById('selectIDTemplateFileButton').addEventListener('click', selectIDTemplateFile);
document.getElementById('selectOutputDirectoryButton').addEventListener('click', selectOutputDirectory);
document.getElementById('generateButton').addEventListener('click', generateIDCards);
