module.exports = {
  branches: ["main"],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
    "@semantic-release/npm",
    "@semantic-release/git",
    [
      "@semantic-release/github",
      { releaseNameTemplate: "webapp/<%= nextRelease.gitTag %>" },
    ],
  ],
  repositoryUrl: "https://github.com/flowck/dobermann.git",
};
