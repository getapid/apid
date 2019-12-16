/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

import React from 'react';
import classnames from 'classnames';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';

const features = [
  {
    title: <>Made for developers</>,
    imageUrl: 'img/illustrations/programming_code_2.svg',
    description: (
      <>
        Write tests in YAML and execute with a single command. Works everywhere,
        even on my machine a developer on a computer
      </>
    )
  },
  {
    title: <>Easily integratable</>,
    imageUrl: 'img/illustrations/data_maintenance.svg',
    description: (
      <>
        APId works with any continuous integration tool, just download the
        binary or use the docker image
      </>
    )
  },
  {
    title: <>Testing framework</>,
    imageUrl: 'img/illustrations/complete_task_1.svg',
    description: (
      <>Run before deployments, make sure your service hasn't regressed</>
    )
  },
  {
    title: <>Versatile</>,
    imageUrl: 'img/illustrations/start_up.svg',
    description: (
      <>
        APId supports shell command execution to provide all the flexibility you
        might need
      </>
    )
  }
];

function Feature({ imageUrl, title, description }) {
  const imgUrl = useBaseUrl(imageUrl);
  return (
    <div className={classnames('col col--3', styles.feature)}>
      {imgUrl && (
        <div className="text--center">
          <img className={styles.featureImage} src={imgUrl} alt={title} />
        </div>
      )}
      <h3 className="text--center">{title}</h3>
      <p>{description}</p>
    </div>
  );
}

function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="API health and performance monitoring"
    >
      <header className={classnames('hero', styles.heroBanner)}>
        <div className="container">
          <div class="row">
            <div class="col col--6">
              <img
                alt="APId sample configuration"
                src={useBaseUrl('img/apid_yaml.png')}
              />
            </div>

            <div class="col col--6">
              <div style={{ display: 'flex', height: '100%' }}>
                <div style={{ margin: 'auto' }}>
                  <h1 className="hero__title">
                    API testing doesn't have to be tedious
                  </h1>
                  <p className="hero__subtitle">
                    Powerful declarative end-to-end testing for APIs that works
                    for you! No coding required. Simple to run on any continuous
                    integration tool.
                  </p>
                  <div className={styles.buttons}>
                    <Link
                      className={classnames(
                        'button button--primary button--lg',
                        styles.download
                      )}
                      to={useBaseUrl('docs/')}
                    >
                      Download
                    </Link>

                    <Link
                      className={classnames(
                        'button button--outline button--secondary button--lg',
                        styles.getStarted
                      )}
                      to={useBaseUrl('docs')}
                    >
                      Get Started
                    </Link>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>
      <main>
        {features && features.length && (
          <section className={styles.features}>
            <div className="container">
              <div className="row">
                {features.map((props, idx) => (
                  <Feature key={idx} {...props} />
                ))}
              </div>
            </div>
          </section>
        )}
      </main>
    </Layout>
  );
}

export default Home;
