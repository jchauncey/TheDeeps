module.exports = {
  moduleNameMapper: {
    '^axios$': 'jest-mock-axios',
    '^react-router-dom$': '<rootDir>/src/__tests__/mocks/react-router-dom.ts',
    '^../../components/CharacterCard$': '<rootDir>/src/__tests__/mocks/CharacterCard.tsx'
  },
  setupFilesAfterEnv: ['<rootDir>/src/setupTests.ts'],
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.(ts|tsx)$': 'ts-jest'
  }
}; 