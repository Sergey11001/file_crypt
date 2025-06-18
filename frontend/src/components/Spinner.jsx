import './Spinner.css';

export default function Spinner({ size = 'medium', color = 'primary' }) {
    const sizeClasses = {
        small: 'spinner-small',
        medium: 'spinner-medium',
        large: 'spinner-large'
    };

    const colorClasses = {
        primary: 'spinner-primary',
        secondary: 'spinner-secondary',
        white: 'spinner-white'
    };

    return (
        <div
            className={`spinner ${sizeClasses[size]} ${colorClasses[color]}`}
            aria-label="Loading"
            role="status"
        />
    );
}