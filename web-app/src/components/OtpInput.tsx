import React, { useRef } from "react";
import '../css/OtpInput.css';

type Props = {
    value: string[];
    setValue: (val: string[]) => void;
};

export default function OtpInput({ value, setValue }: Props) {
    const inputs = useRef<(HTMLInputElement | null)[]>([]);

    const handleChange = (text: string, index: number) => {
        const newValue = [...value];
        const sanitizedText = text.replace(/[^0-9]/g, "").slice(-1);
        newValue[index] = sanitizedText;
        setValue(newValue);

        if (sanitizedText && index < value.length - 1) {
            inputs.current[index + 1]?.focus();
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent, index: number) => {
        if (e.key === 'Backspace') {
            if (!value[index] && index > 0) {
                inputs.current[index - 1]?.focus();
            }
        }
    };

    return (
        <div className="otp-row">
            {value.map((digit, i) => (
                <input
                    key={i}
                    ref={(el) => { inputs.current[i] = el; }}
                    className="otp-box"
                    value={digit}
                    onChange={(e) => handleChange(e.target.value, i)}
                    onKeyDown={(e) => handleKeyDown(e, i)}
                    type="text"
                    inputMode="numeric"
                    maxLength={1}
                    autoFocus={i === 0}
                />
            ))}
        </div>
    );
}