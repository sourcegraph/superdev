import React from 'react';
import { render, screen } from '@testing-library/react';
import { processEscapeCodes } from '../utils/escapeCodeHandler';

describe('escapeCodeHandler', () => {
  test('converts ANSI escape codes to HTML', () => {
    const input = '\u001b[31mThis is red text\u001b[0m';
    const result = processEscapeCodes(input);
    
    // Should contain HTML for colored text
    expect(result).toContain('<span');
    expect(result).toContain('color:');
  });

  test('handles empty input', () => {
    const result = processEscapeCodes('');
    expect(result).toBe('');
  });

  test('passes through regular text', () => {
    const input = 'Hello World';
    const result = processEscapeCodes(input);
    expect(result).toBe('Hello World');
  });

  test('handles multiple colors and styles', () => {
    // Add more complex ANSI codes to ensure we get at least 3 spans
    const input = '\u001b[31mRed\u001b[0m \u001b[32mGreen\u001b[0m \u001b[1;34mBold Blue\u001b[0m';
    const result = processEscapeCodes(input);
    

    // Should contain multiple spans
    expect(result.match(/<span/g).length).toBeGreaterThanOrEqual(3);
  });
});