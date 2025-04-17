import React from 'react';
import './Header.css';
import ampsterLogo from '../../assets/ampster.png';

const Header: React.FC = () => {
  return (
    <header className="header">
      <div className="logo-container">
        {/* Component 1: SuperDev brand logo */}
        <img src={ampsterLogo} alt="SuperDev Logo" className="logo" />
        <span className="brand-name">SuperDev</span>
      </div>
    </header>
  );
};

export default Header;