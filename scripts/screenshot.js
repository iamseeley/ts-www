const puppeteer = require('puppeteer');
const fs = require('fs');
const path = require('path');

async function takeScreenshot(tempHTMLFilePath, outputPath) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();

    // Convert file path to a URL
    const fileUrl = 'file://' + path.resolve(tempHTMLFilePath);

    await page.goto(fileUrl, { waitUntil: 'networkidle0' });
    await page.screenshot({ path: outputPath, fullPage: true });
    await browser.close();
}

const tempHTMLFilePath = process.argv[2];
const outputPath = process.argv[3];

takeScreenshot(tempHTMLFilePath, outputPath);
