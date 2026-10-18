import { useState, useEffect } from 'react';

type Theme = 'system' | 'light' | 'dark';

const order: Theme[] = ['system', 'light', 'dark'];

const names: Record<Theme, string> = {
  system: 'the system theme',
  light: 'the light theme',
  dark: 'the dark theme',
};

/**
 * ThemeSelect — React Island for choosing the theme
 *
 * Features:
 * - Cycles through system, light and dark; system is the default
 * - Persists a light or dark choice in localStorage, and forgets it for system
 * - Follows the system setting as it changes, while system is chosen
 * - No FOUC (initial state set by inline script in BaseLayout)
 */
export default function ThemeSelect() {
  const [theme, setTheme] = useState<Theme>('system');

  useEffect(() => {
    // Read the actual state from the DOM (set by BaseLayout inline script)
    const chosen = document.documentElement.dataset.theme;
    setTheme(chosen === 'light' || chosen === 'dark' ? chosen : 'system');
  }, []);

  function cycle() {
    const next = order[(order.indexOf(theme) + 1) % order.length];
    setTheme(next);
    const root = document.documentElement;
    root.dataset.theme = next;
    root.classList.toggle(
      'dark',
      next === 'dark' || (next === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches),
    );
    try {
      if (next === 'system') {
        localStorage.removeItem('theme');
      } else {
        localStorage.setItem('theme', next);
      }
    } catch {
      // Storage is blocked: the choice lasts until the page is left.
    }
  }

  const next = order[(order.indexOf(theme) + 1) % order.length];
  const label = `Using ${names[theme]}. Switch to ${names[next]}`;
  const icon = (shown: Theme) =>
    `h-5 w-5 transition-all duration-300 ${
      theme === shown ? 'scale-100 rotate-0 opacity-100' : 'scale-0 rotate-90 opacity-0'
    } absolute`;

  return (
    <button
      onClick={cycle}
      aria-label={label}
      title={label}
      className="cursor-pointer relative flex h-9 w-9 items-center justify-center rounded-lg text-[var(--text-secondary)] transition-colors hover:bg-[var(--bg-secondary)] hover:text-[var(--text-primary)]"
    >
      {/* Monitor icon */}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        className={icon('system')}
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
        />
      </svg>
      {/* Sun icon */}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        className={icon('light')}
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
        />
      </svg>
      {/* Moon icon */}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        className={icon('dark')}
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
        />
      </svg>
    </button>
  );
}
