import Convert from 'ansi-to-html';
import DOMPurify from 'dompurify';

// Initialize the converter with options
const convert = new Convert({
  fg: '#000',
  bg: '#fff',
  newline: true,
  escapeXML: true,
  stream: false,
  colors: {
    0: '#000',    // Black
    1: '#e00',    // Red
    2: '#0a0',    // Green
    3: '#aa0',    // Yellow
    4: '#00a',    // Blue
    5: '#a0a',    // Magenta
    6: '#0aa',    // Cyan
    7: '#aaa',    // White
  }
});

/**
 * Process terminal escape codes in text and convert to safe HTML
 * @param {string} text - The input text containing ANSI escape codes
 * @returns {string} Sanitized HTML with escape codes converted to styled spans
 */
export function processEscapeCodes(text) {
  if (!text) return '';
  
  try {
    // Make sure reset codes are properly separated to create distinct spans
    const processedText = text.replace(/\u001b\[0m\s*\u001b\[/g, '\u001b[0m \u001b[');
    
    // Convert ANSI escape codes to HTML
    const html = convert.toHtml(processedText);
    
    // Sanitize the HTML to prevent XSS attacks
    return DOMPurify.sanitize(html);
  } catch (error) {
    console.error('Error processing escape codes:', error);
    return text; // Return original text on error
  }
}