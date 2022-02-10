module.exports = {
    content: ["./templates/**/*.{html,js}"],
    theme: {
        fontFamily: {
            body: ['Raleway'],
            display: ['Prompt']
        },
        extend: {
            colors: {
                purple: {
                    100: '#9D4EDD',
                    200: '#7B2CBF',
                    300: '#5A189A',
                    400: '#3C096C',
                    500: '#240046'
                },
                orange: {
                    100: '#FF9E00',
                    200: '#FF9100',
                    300: '#FF8500',
                    400: '#FF7900',
                    500: '#FF6D00',
                },
            },
        },
    },
    plugins: [],
}