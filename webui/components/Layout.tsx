import Head from 'next/head';

export const Layout = ({ title, children }) => {
  return (
    <>
      <Head>
        <title>{title}</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {children}
    </>
  );
};
